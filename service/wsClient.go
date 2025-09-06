package service

import (
	"comfyui_endpoint/client"
	"comfyui_endpoint/logger"
	"comfyui_endpoint/model"
	"comfyui_endpoint/ws"
	"encoding/base64"
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/tidwall/gjson"
)

var (
	clientMap    = make(map[string]WsClient)
	clientsMutex = sync.Mutex{}
)

// 限制同个uuid下一个ws连接

type WsClient interface {
	Start() error
	Close() error
	SetClient(client ws.WsClient)
}

type wsClient struct {
	addr   string
	uuid   string
	client ws.WsClient
}

func NewWsClient(uuid, addr string) WsClient {
	if clientMap[uuid] == nil {
		client := wsClient{
			addr: addr,
			uuid: uuid,
		}
		clientMap[uuid] = &client
		return clientMap[uuid]
	} else {
		return clientMap[uuid]
	}

}

func (w *wsClient) SetClient(client ws.WsClient) {
	w.client = client
}

func (w *wsClient) Start() error {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()
	var err error
	fmt.Println(clientMap[w.uuid])
	if w.client == nil {
		wsclient, err := ws.NewWsClient(w.addr,
			func(msg string) {
				// TODO ws消息业务逻辑
				if gjson.Get(msg, "type").String() == "progress_state" {
					{
						promptId := gjson.Get(msg, "data.prompt_id").String()
						resp, err := client.RestyClient.R().Get(fmt.Sprintf("http://%s/history/%s", w.addr, promptId))
						if err != nil {
							logger.Error("ws获取进度失败: %s", err.Error())
						}
						respString := resp.String()
						if gjson.Get(respString, promptId).String() != "" {
							filename := gjson.Get(respString, promptId+".outputs.*.images.0.filename").String()

							if filename != "" {
								{
									// 同步策略
									// 将结果写入 channel
									syncPromptMapMutex.RLock() // 使用读锁
									if ch, ok := syncPromptMap[promptId]; ok {
										select {
										case ch <- struct{ FileName string }{FileName: filename}:
										default:
											logger.Debug("无法发送到 channel，可能已关闭或已满: %s", promptId)
										}
									}
									syncPromptMapMutex.RUnlock()
								}

								{

									// 异步策略
									asyncPromptMapMutex.RLock()
									ch, ok := asyncPromptMap[promptId]
									if !ok {
										asyncPromptMapMutex.RUnlock()
										return
									}
									res := <-ch
									if res.CallbackUrl != "" {
										// 处理接收到的数据
										imgResp, err := client.RestyClient.R().Get(fmt.Sprintf("http://%s/view?filename=%s", res.Addr, filename))
										body := make(map[string]interface{})
										if err != nil {
											body = map[string]interface{}{
												"code":    500,
												"message": fmt.Sprintf("获取图片失败: %s", err.Error()),
												"data":    "",
											}
										} else {
											// 将图片数据转换为base64字符串
											imageBase64 := base64.StdEncoding.EncodeToString(imgResp.Body())
											body = map[string]interface{}{
												"code":    200,
												"message": "回调成功",
												"data": map[string]interface{}{
													"file_name":          filename,
													"prompt_id":          promptId,
													"uri":                fmt.Sprintf("http://%s/view?filename=%s", res.Addr, filename),
													"uid":                w.uuid,
													"binary_data_base64": []string{imageBase64},
												},
											}
										}

										// 发送回调请求
										client.RestyClient.R().SetDebug(true).SetBody(body).Post(res.CallbackUrl)
									}

									asyncPromptMapMutex.RUnlock()
									// 删除用完的ch
									asyncPromptMapMutex.Lock()
									delete(asyncPromptMap, promptId)
									close(ch)
									asyncPromptMapMutex.Unlock()
								}

							}
						}
					}
				}
			},
			func(err error) {
				if _, ok := err.(*websocket.CloseError); ok {
					clientsMutex.Lock()
					defer clientsMutex.Unlock()
					// 数据库标记ws close状态
					SqliteDb().Model(&model.ComfyAppInfo{}).Where("uid = ?", w.uuid).Update("ws_status", "close")
				}
			},
			func(client ws.WsClient) {
				clientsMutex.Lock()
				defer clientsMutex.Unlock()
				clientMap[w.uuid].SetClient(client)
				// 数据库标记ws connect状态
				SqliteDb().Model(&model.ComfyAppInfo{}).Where("uid = ?", w.uuid).Update("ws_status", "connect")
			},
		)

		if err != nil {
			logger.Error("ws client start error %s", err.Error())
		}
		w.SetClient(wsclient)
		clientMap[w.uuid] = w
	}

	return err
}

func (w *wsClient) Close() error {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()
	if w.client != nil {
		err := w.client.Close()
		if err != nil {
			return err
		}
		delete(clientMap, w.uuid)
	}
	return nil
}

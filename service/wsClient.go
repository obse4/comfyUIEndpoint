package service

import (
	"comfyui_endpoint/client"
	"comfyui_endpoint/logger"
	"comfyui_endpoint/model"
	"comfyui_endpoint/ws"
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
					promptId := gjson.Get(msg, "data.prompt_id").String()
					resp, err := client.RestyClient.R().Get(fmt.Sprintf("http://%s/history/%s", w.addr, promptId))
					if err != nil {
						logger.Error("ws获取进度失败: %s", err.Error())
					}
					respString := resp.String()
					if gjson.Get(respString, promptId).String() != "" {
						filename := gjson.Get(respString, promptId+".outputs.*.images.0.filename").String()

						if filename != "" {
							// 将结果写入 channel
							if ch, ok := syncPromptMap[promptId]; ok {
								ch <- struct{ FileName string }{FileName: filename}
								//? close(ch) // 关闭 channel 防止重复写入(因ws会发送两次相同prompt,此处不可关闭，会导致panic) 内存 待优化
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

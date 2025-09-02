package ws

import (
	"comfyui_endpoint/logger"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

type WsClient interface {
	Close() error
}

type wsClient struct {
	url  string
	conn *websocket.Conn
}

func (s *wsClient) Close() error {
	return s.conn.Close()
}

func reconnect(uri string) (*websocket.Conn, *http.Response, error) {
	return websocket.DefaultDialer.Dial(uri, nil)
}

func NewWsClient(addr string, f func(msg string), dealErr func(err error), dealReconnect func(client WsClient)) (ws WsClient, err error) {
	client := &wsClient{url: strings.Join([]string{"ws:/", addr, "ws"}, "/")}
	client.conn, _, err = reconnect(client.url)
	if err != nil {
		logger.Error("ws 连接失败: %s", err.Error())
		return
	}

	logger.Info("ws 已连接到服务器 %s", client.url)

	// 设置中断信号监听，优雅关闭
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// 消息接收通道
	messageChan := make(chan []byte)
	// 错误通道
	errorChan := make(chan error)

	// 启动goroutine接收消息
	go func() {
		for {
			_, message, err := client.conn.ReadMessage()
			if err != nil {
				errorChan <- err
			ReConnectWS:
				for {
					client.conn, _, err = reconnect(client.url)
					logger.Debug("尝试重新连接ws")
					if err != nil {
						time.Sleep(time.Second * 5)
					} else {
						dealReconnect(client.conn)
						break ReConnectWS
					}
				}

			}
			messageChan <- message
		}
	}()

	go func() {
		// 主循环处理消息和事件
		for {
			select {
			case <-interrupt:
				logger.Debug("接收到中断信号，关闭连接...")

				// 发送关闭消息
				err = client.conn.WriteMessage(websocket.CloseMessage,
					websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
				if err != nil {
					logger.Debug("发送关闭消息错误: %s", err)
				}
				return

			case err = <-errorChan:
				logger.Error("发生错误: %s", err)
				dealErr(err)
			case message := <-messageChan:
				logger.Info("接收到消息: %s", string(message))
				f(string(message))
			}
		}
	}()

	return client, err
}

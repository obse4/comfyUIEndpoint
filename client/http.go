package client

import (
	"net"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
)

var RestyClient *resty.Client

func InitRestyClient() {

	RestyClient = resty.New().SetTransport(&http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   60 * time.Second,
			KeepAlive: 60 * time.Second,
		}).DialContext,
		DisableKeepAlives:     true, // 禁用长链接
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          1000,             // 允许的最大空闲连接数
		MaxIdleConnsPerHost:   2000,             // 每个主机允许的最大空闲连接数 接池中为每个主机维护的最大空闲连接数。
		IdleConnTimeout:       90 * time.Second, // 空闲连接的超时时间。
		TLSHandshakeTimeout:   15 * time.Second, // 握手的超时时间。
		ExpectContinueTimeout: 1 * time.Second,
	}).SetTimeout(10*60*time.Second).
		SetHeader("User-Agent", "").
		SetDebug(false)
}

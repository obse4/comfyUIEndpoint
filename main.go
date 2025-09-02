package main

import (
	"comfyui_endpoint/client"
	"comfyui_endpoint/config"
	"comfyui_endpoint/controller"
	"comfyui_endpoint/database"
	_ "comfyui_endpoint/docs"
	"comfyui_endpoint/logger"
	"comfyui_endpoint/service"
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	// 初始化配置项
	config.InitConfig()
	// 初始化全局http客户端
	client.InitRestyClient()
	// 初始化自动化日志收集
	logger.InitLogger(&logger.LogConfig{LogLevel: logger.LogLeveL(config.Global.Log.Level)})
	database.InitSqlite(&config.Global.Sqlite)
	service.NewComfyAppService().InitWs()

	// 获取当前cpu数量
	maxCpuNum := runtime.NumCPU()
	// 最大cpu启动
	runtime.GOMAXPROCS(maxCpuNum)

	// 配置gin mode
	gin.SetMode(config.Global.HttpServer.Mode)

	router := gin.New()
	service.GinRouter = router

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, func(c *ginSwagger.Config) { c.PersistAuthorization = true }))

	{
		router.GET("health", func(ctx *gin.Context) {
			logger.Debug("health check success")
			ctx.String(200, "ok")
		})
	}

	{ // endpointParam
		router.POST("endpointParam/set", controller.EndpointParamSet)
		router.POST("endpointParam/find", controller.EndpointParamFid)
	}

	{ // comfyApp
		router.POST("comfyApp/create", controller.ComfyAppCreate)
		router.POST("comfyApp/remove", controller.ComfyAppRemove)
		router.POST("comfyApp/update", controller.ComfyAppUpdate)
		router.POST("comfyApp/index", controller.ComfyAppIndex)
		router.POST("comfyApp/wsRestart", controller.ComfyAppRestartWs)

	}

	{ // endpoint
		router.POST("endpoint/create", controller.EndpointCreate)
		router.POST("endpoint/update", controller.EndpointUpdate)
		router.POST("endpoint/index", controller.EndpointIndex)

	}

	service.InitSyncHandle(router)

	var port string
	flag.StringVar(&port, "p", "9518", "server port")

	if !strings.Contains(port, ":") {
		port = fmt.Sprintf(":%s", port)
	}
	logger.Info("server start at %s", port)
	srv := &http.Server{
		Addr:    port,
		Handler: router,
	}

	go func() {
		// 开启goroutine启动服务
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Server Fatal: %v", err)
		}
	}()

	// 等待中断信号来优雅地关闭服务器，为关闭服务器操作设置一个5秒的超时
	quit := make(chan os.Signal, 1)
	// kill 默认会发送 syscall.SIGTERM 信号
	// kill -2 发送 syscall.SIGINT 信号，我们常用的Ctrl+C就是触发系统SIGINT信号
	// kill -9 发送 syscall.SIGKILL 信号，但是不能被捕获，所以不需要添加它
	// signal.Notify把收到的 syscall.SIGINT或syscall.SIGTERM 信号转发给quit
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	// 阻塞在此，当接收到上述两种信号时才会往下执行
	<-quit

	logger.Info("--- Shutting Down Server ---")
	// 创建一个10秒超时的context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()
	// 10秒内优雅关闭服务（将未处理完的请求处理完再关闭服务），超过10秒就超时退出
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server Forced To Shutdown: %v", err)
	}

	logger.Info("--- Server Exiting ---")
	defer logger.Info("--- Server Closed ---")
}

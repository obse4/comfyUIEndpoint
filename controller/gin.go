package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 接口错误返回函数
func GinError(ctx *gin.Context, err error, message string, code int) (isErr bool) {
	isErr = false
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    code,
			"message": fmt.Sprintf("%s [%s]", message, err.Error()),
			"data":    err.Error(),
		})

		isErr = true
	}
	return
}

// 接口返回
func GinReply(ctx *gin.Context, message string, code int, data interface{}) {
	ctx.JSON(http.StatusOK, gin.H{
		"code":    code,
		"message": message,
		"data":    data,
	})
}

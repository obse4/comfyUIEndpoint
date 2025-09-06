package service

import "github.com/gin-gonic/gin"

var GinRouter *gin.Engine

func GetRouter() *gin.Engine {
	return GinRouter
}

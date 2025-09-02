package controller

import (
	"comfyui_endpoint/dto/request"
	"comfyui_endpoint/model"
	"comfyui_endpoint/service"

	"github.com/gin-gonic/gin"
)

var endpointParamService = service.NewEndpointParamService()

// @Tags EndpointParam
// @Summary 端点接口参数 配置
// @Description 端点接口参数 配置
// @Accept json
// @Produce json
// @Param request body request.EndpointParamSetRequest true "请求参数"
// @Success 200 {object} response.CommonResponse "返回值"
// @Router /endpointParam/set [post]
func EndpointParamSet(ctx *gin.Context) {
	var params request.EndpointParamSetRequest

	err := ctx.ShouldBindJSON(&params)

	if GinError(ctx, err, "参数错误", 400) {
		return
	}

	err = endpointParamService.Set(params)

	if GinError(ctx, err, "创建失败", 500) {
		return
	}

	GinReply(ctx, "创建成功", 200, nil)
}

// @Tags EndpointParam
// @Summary 端点接口参数 查询
// @Description 端点接口参数 查询
// @Accept json
// @Produce json
// @Param request body request.EndpointParamFindRequest true "请求参数"
// @Success 200 {object} response.CommonResponse{data=[]model.EndpointParam} "返回值"
// @Router /endpointParam/find [post]
func EndpointParamFid(ctx *gin.Context) {
	var params request.EndpointParamFindRequest

	err := ctx.ShouldBindJSON(&params)

	if GinError(ctx, err, "参数错误", 400) {
		return
	}

	var data []model.EndpointParam

	data, err = endpointParamService.Find(params)

	if GinError(ctx, err, "查询失败", 500) {
		return
	}

	GinReply(ctx, "查询成功", 200, data)
}

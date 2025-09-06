package controller

import (
	"comfyui_endpoint/dto/request"
	"comfyui_endpoint/dto/response"
	"comfyui_endpoint/service"

	"github.com/gin-gonic/gin"
)

var endpointService = service.NewEndpointService()

// @Tags Endpoint
// @Summary 端点 创建
// @Description 端点 创建
// @Accept json
// @Produce json
// @Param request body request.EndpointCreateRequest true "请求参数"
// @Success 200 {object} response.CommonResponse "返回值"
// @Router /endpoint/create [post]
func EndpointCreate(ctx *gin.Context) {
	var params request.EndpointCreateRequest

	err := ctx.ShouldBindJSON(&params)

	if GinError(ctx, err, "参数错误", 400) {
		return
	}

	err = endpointService.Create(params)

	if GinError(ctx, err, "创建失败", 500) {
		return
	}

	GinReply(ctx, "创建成功", 200, nil)
}

// @Tags Endpoint
// @Summary 端点 更新
// @Description 端点 更新
// @Accept json
// @Produce json
// @Param request body request.EndpointUpdateRequest true "请求参数"
// @Success 200 {object} response.CommonResponse "返回值"
// @Router /endpoint/update [post]
func EndpointUpdate(ctx *gin.Context) {
	var params request.EndpointUpdateRequest

	err := ctx.ShouldBindJSON(&params)

	if GinError(ctx, err, "参数错误", 400) {
		return
	}

	err = endpointService.Update(params)

	if GinError(ctx, err, "更新失败", 500) {
		return
	}

	GinReply(ctx, "更新成功", 200, nil)
}

// @Tags Endpoint
// @Summary 端点 分页查询
// @Description 端点 分页查询
// @Accept json
// @Produce json
// @Param request body request.EndpointIndexRequest true "请求参数"
// @Success 200 {object} response.CommonResponse{data=response.EndpointIndexResponse} "返回值"
// @Router /endpoint/index [post]
func EndpointIndex(ctx *gin.Context) {
	var params request.EndpointIndexRequest

	err := ctx.ShouldBindJSON(&params)

	if GinError(ctx, err, "参数错误", 400) {
		return
	}

	var data response.EndpointIndexResponse
	data, err = endpointService.Index(params)

	if GinError(ctx, err, "查询失败", 500) {
		return
	}

	GinReply(ctx, "查询成功", 200, data)
}

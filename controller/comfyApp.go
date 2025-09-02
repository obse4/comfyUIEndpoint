package controller

import (
	"comfyui_endpoint/dto/request"
	"comfyui_endpoint/dto/response"
	"comfyui_endpoint/service"

	"github.com/gin-gonic/gin"
)

var comfyAppService = service.NewComfyAppService()

// @Tags ComfyApp
// @Summary comfy应用 创建
// @Description comfy应用 创建
// @Accept json
// @Produce json
// @Param request body request.ComfyAppCreateRequest true "请求参数"
// @Success 200 {object} response.CommonResponse "返回值"
// @Router /comfyApp/create [post]
func ComfyAppCreate(ctx *gin.Context) {
	var params request.ComfyAppCreateRequest

	err := ctx.ShouldBindJSON(&params)

	if GinError(ctx, err, "参数错误", 400) {
		return
	}

	err = comfyAppService.Create(params)

	if GinError(ctx, err, "创建失败", 500) {
		return
	}

	GinReply(ctx, "创建成功", 200, nil)
}

// @Tags ComfyApp
// @Summary comfy应用 删除
// @Description comfy应用 删除
// @Accept json
// @Produce json
// @Param request body request.ComfyAppRemoveRequest true "请求参数"
// @Success 200 {object} response.CommonResponse "返回值"
// @Router /comfyApp/remove [post]
func ComfyAppRemove(ctx *gin.Context) {
	var params request.ComfyAppRemoveRequest

	err := ctx.ShouldBindJSON(&params)

	if GinError(ctx, err, "参数错误", 400) {
		return
	}

	err = comfyAppService.Remove(params)

	if GinError(ctx, err, "删除失败", 500) {
		return
	}

	GinReply(ctx, "删除成功", 200, nil)
}

// @Tags ComfyApp
// @Summary comfy应用 更新
// @Description comfy应用 更新
// @Accept json
// @Produce json
// @Param request body request.ComfyAppUpdateRequest true "请求参数"
// @Success 200 {object} response.CommonResponse "返回值"
// @Router /comfyApp/update [post]
func ComfyAppUpdate(ctx *gin.Context) {
	var params request.ComfyAppUpdateRequest

	err := ctx.ShouldBindJSON(&params)

	if GinError(ctx, err, "参数错误", 400) {
		return
	}

	err = comfyAppService.Update(params)

	if GinError(ctx, err, "更新失败", 500) {
		return
	}

	GinReply(ctx, "更新成功", 200, nil)
}

// @Tags ComfyApp
// @Summary comfy应用 分页查询
// @Description comfy应用 分页查询
// @Accept json
// @Produce json
// @Param request body request.ComfyAppIndexRequest true "请求参数"
// @Success 200 {object} response.CommonResponse{data=response.ComfyAppIndexResponse} "返回值"
// @Router /comfyApp/index [post]
func ComfyAppIndex(ctx *gin.Context) {
	var params request.ComfyAppIndexRequest

	err := ctx.ShouldBindJSON(&params)

	if GinError(ctx, err, "参数错误", 400) {
		return
	}

	var data response.ComfyAppIndexResponse
	data, err = comfyAppService.Index(params)

	if GinError(ctx, err, "查询失败", 500) {
		return
	}

	GinReply(ctx, "查询成功", 200, data)
}

// @Tags ComfyApp
// @Summary comfy应用 ws重启
// @Description comfy应用 ws重启
// @Accept json
// @Produce json
// @Param request body request.ComfyAppRestartWsRequest true "请求参数"
// @Success 200 {object} response.CommonResponse "返回值"
// @Router /comfyApp/wsRestart [post]
func ComfyAppRestartWs(ctx *gin.Context) {
	var params request.ComfyAppRestartWsRequest

	err := ctx.ShouldBindJSON(&params)

	if GinError(ctx, err, "参数错误", 400) {
		return
	}

	err = comfyAppService.RestartWs(params)

	if GinError(ctx, err, "ws重启失败", 500) {
		return
	}

	GinReply(ctx, "ws重启成功", 200, nil)
}

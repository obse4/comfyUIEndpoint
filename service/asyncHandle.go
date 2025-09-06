package service

import (
	"comfyui_endpoint/client"
	"comfyui_endpoint/model"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/sjson"
)

var asyncPromptMap = make(map[string]chan struct {
	CallbackUrl string
	Addr        string
})
var asyncPromptMapMutex = sync.RWMutex{} // 添加读写锁

func InitAsyncHandle(r *gin.Engine) {
	var list []model.Endpoint
	SqliteDb().Model(&model.Endpoint{}).Find(&list)

	for _, v := range list {
		RegisterAsyncHandle(r, v.Path)
	}
}

func RegisterAsyncHandle(r *gin.Engine, path string) {
	r.POST(path, AsyncHandle)
}

func AsyncHandle(ctx *gin.Context) {
	path := strings.TrimPrefix(ctx.FullPath(), "/")

	var params = make(map[string]interface{})
	err := ctx.ShouldBindJSON(&params)

	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": fmt.Sprintf("%s [%s]", "参数错误", err.Error()),
			"data":    err.Error(),
		})
		return
	}
	var endpoint model.Endpoint

	SqliteDb().Model(&model.Endpoint{}).Where("path = ?", path).First(&endpoint)

	if endpoint.Id == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": fmt.Sprintf("%s [%s]", "找不到路径对应endpoint", path),
			"data":    "",
		})
		return
	}

	var comfyApp model.ComfyApp

	uid, ok := params["uid"].(string)
	fmt.Println(uid)
	if uid == "" && !ok {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "uid不能为空",
			"data":    "",
		})
		return

	}
	SqliteDb().Model(&model.ComfyApp{}).Where("uid = ?", uid).First(&comfyApp)
	if comfyApp.Id == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": fmt.Sprintf("%s [%s]", "找不到uid对应app", uid),
			"data":    "",
		})
		return
	}

	for k, v := range params {
		switch k {
		case "uid":
		case "callback_url":
			callbackUrl, ok := v.(string)

			if !ok && endpoint.CallbackUrl == "" {
				ctx.JSON(http.StatusOK, gin.H{
					"code":    500,
					"message": fmt.Sprintf("%s [%s]", "找不到对应callbackUrl", path),
					"data":    "",
				})
				return
			}
			if callbackUrl != "" {
				endpoint.CallbackUrl = callbackUrl
			}
		default:
			// 其他参数用于修改json文件
			endpointParam, err := NewEndpointParamService().FindOne(path, k)
			if err != nil {
				ctx.JSON(http.StatusOK, gin.H{
					"code":    500,
					"message": fmt.Sprintf("%s [%s]", "参数错误", err.Error()),
					"data":    err.Error(),
				})
				return
			}
			endpoint.ApiJson, err = sjson.Set(endpoint.ApiJson, endpointParam.JsonKey, v)

			if err != nil {
				ctx.JSON(http.StatusOK, gin.H{
					"code":    500,
					"message": fmt.Sprintf("%s [%s]", "参数写入错误", err.Error()),
					"data":    err.Error(),
				})
				return
			}
		}
	}

	if endpoint.CallbackUrl == "" {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": fmt.Sprintf("%s [%s]", "找不到对应callbackUrl", path),
			"data":    "",
		})
		return
	}

	var apiJsonMap = make(map[string]interface{})
	json.Unmarshal([]byte(endpoint.ApiJson), &apiJsonMap)

	promptResp, err := client.RestyClient.R().SetBody(map[string]interface{}{
		"prompt": apiJsonMap,
	}).Post(fmt.Sprintf("http://%s/api/prompt", comfyApp.Addr))

	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": fmt.Sprintf("%s [%s]", "prompt请求错误", err.Error()),
			"data":    err.Error(),
		})
		return
	}
	var promptRespData struct {
		PromptId   string      `json:"prompt_id"`
		NodeErrors interface{} `json:"node_errors"`
	}
	promptRespBody := promptResp.Body()
	json.Unmarshal(promptRespBody, &promptRespData)

	if promptRespData.PromptId == "" {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": fmt.Sprintf("%s [%v]", "prompt请求错误", promptRespData.NodeErrors),
			"data":    promptResp.String(),
		})
		return
	}
	asyncPromptMapMutex.Lock()
	asyncPromptMap[promptRespData.PromptId] = make(chan struct {
		CallbackUrl string
		Addr        string
	}, 1)
	asyncPromptMap[promptRespData.PromptId] <- struct {
		CallbackUrl string
		Addr        string
	}{CallbackUrl: endpoint.CallbackUrl, Addr: comfyApp.Addr}
	asyncPromptMapMutex.Unlock()

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "任务创建成功",
		"data": map[string]interface{}{
			"prompt_id": promptRespData.PromptId,
			"uid":       comfyApp.Uid,
		},
	})
}

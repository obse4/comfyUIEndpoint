package service

import (
	"comfyui_endpoint/client"
	"comfyui_endpoint/model"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/sjson"
)

var GinRouter *gin.Engine
var syncPromptMap = make(map[string]chan struct{ FileName string })
var syncPromptMapMutex = sync.RWMutex{} // 添加读写锁
func GetRouter() *gin.Engine {
	return GinRouter
}

// 服务启动时初始化
func InitSyncHandle(r *gin.Engine) {
	var list []model.Endpoint
	SqliteDb().Model(&model.Endpoint{}).Find(&list)

	for _, v := range list {
		RegisterSyncHandle(r, v.SyncPath)
	}
}

func RegisterSyncHandle(r *gin.Engine, path string) {
	r.POST(path, SyncHandle)
}

func SyncHandle(ctx *gin.Context) {
	// 等待ws告知执行完成
	// 获取图片文件二进制数据
	// 返回二进制数据
	path := strings.TrimPrefix(ctx.FullPath(), "/")
	path = strings.TrimSuffix(path, "/sync")
	// 通过path找到endpoint
	// 拼接参数
	// 调用comfy 获取prompt_id
	// 等待数据
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

	for k, v := range params {
		switch k {
		case "uid":
			if v == "" {
				ctx.JSON(http.StatusOK, gin.H{
					"code":    500,
					"message": "uid不能为空",
					"data":    "",
				})
				return

			}
			SqliteDb().Model(&model.ComfyApp{}).Where("uid = ?", v).First(&comfyApp)
			if comfyApp.Id == 0 {
				ctx.JSON(http.StatusOK, gin.H{
					"code":    500,
					"message": fmt.Sprintf("%s [%s]", "找不到uid对应app", v),
					"data":    "",
				})
				return
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

	var apiJsonMap = make(map[string]interface{})
	json.Unmarshal([]byte(endpoint.ApiJson), &apiJsonMap)
	fmt.Println(apiJsonMap)

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
	syncPromptMapMutex.Lock()
	resultChan := make(chan struct {
		FileName string
	}, 1)
	syncPromptMap[promptRespData.PromptId] = resultChan
	syncPromptMapMutex.Unlock()
	// 等待数据
	var res struct {
		FileName string
	}
	select {
	case res = <-resultChan:
		// 正常处理
		syncPromptMapMutex.Lock()
		delete(syncPromptMap, promptRespData.PromptId)
		close(resultChan)
		syncPromptMapMutex.Unlock()
	case <-time.After(60 * time.Second):
		syncPromptMapMutex.Lock()
		delete(syncPromptMap, promptRespData.PromptId)
		close(resultChan)
		syncPromptMapMutex.Unlock()
		ctx.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "等待超时",
			"data":    "",
		})
		return
	}

	// 处理接收到的数据
	imgResp, err := client.RestyClient.R().Get(fmt.Sprintf("http://%s/view?filename=%s", comfyApp.Addr, res.FileName))

	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": fmt.Sprintf("获取图片失败: %s", err.Error()),
			"data":    "",
		})
		return
	}

	// 将图片数据转换为base64字符串
	imageBase64 := base64.StdEncoding.EncodeToString(imgResp.Body())

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取图片成功",
		"data": map[string]interface{}{
			"file_name":          res.FileName,
			"prompt_id":          promptRespData.PromptId,
			"uri":                fmt.Sprintf("http://%s/view?filename=%s", comfyApp.Addr, res.FileName),
			"uid":                comfyApp.Uid,
			"binary_data_base64": []string{imageBase64},
		},
	})
}

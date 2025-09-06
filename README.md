# ComfyUI Endpoint (comfyUI 工作流自动化端点)
本项目为comfyUI工作流接口化调用服务，旨在实现comfy工作流端点接入与配置

## 快速开始
1. golang环境
2. 运行 `go mod tidy` 
3. 运行 `go run main.go`
4. 访问 `http://localhost:9518/swagger/index.html#/` 

## 配置
1. 参考`config/config.yaml`
2. 配置端口 `comfyui_endpoint -p 8080`
3. 启动服务 `comfyui_endpoint`
4. 访问swagger `http://localhost:9518/swagger/index.html#/`

## 操作
1. 创建comfy应用 `http://localhost:9518/swagger/index.html#/ComfyApp/post_comfyApp_create`

```HTTP
POST /comfyApp/create HTTP/1.1
Host: 127.0.0.1:9518
Content-Type: application/json
Content-Length: 60

{
  "addr": "127.0.0.1:8000",
  "description": "本地comfyUI"
}
```

**注意：** addr只是填 ip+端口


2. 创建端口 `http://localhost:9518/swagger/index.html#/Endpoint/post_endpoint_create`
```HTTP
POST /endpoint/create HTTP/1.1
Host: 127.0.0.1:9518
Content-Type: application/json
Content-Length: 1374

{
  "api_json": "{\"3\":{\"inputs\":{\"seed\":[\"10\",0],\"steps\":20,\"cfg\":8,\"sampler_name\":\"euler\",\"scheduler\":\"normal\",\"denoise\":1,\"model\":[\"4\",0],\"positive\":[\"6\",0],\"negative\":[\"7\",0],\"latent_image\":[\"5\",0]},\"class_type\":\"KSampler\",\"_meta\":{\"title\":\"K采样器\"}},\"4\":{\"inputs\":{\"ckpt_name\":\"v1-5-pruned-emaonly-fp16.safetensors\"},\"class_type\":\"CheckpointLoaderSimple\",\"_meta\":{\"title\":\"Checkpoint加载器（简易）\"}},\"5\":{\"inputs\":{\"width\":512,\"height\":512,\"batch_size\":1},\"class_type\":\"EmptyLatentImage\",\"_meta\":{\"title\":\"空Latent图像\"}},\"6\":{\"inputs\":{\"text\":\"beautiful scenery nature glass bottle landscape, , purple galaxy bottle,\",\"clip\":[\"4\",1]},\"class_type\":\"CLIPTextEncode\",\"_meta\":{\"title\":\"CLIP文本编码\"}},\"7\":{\"inputs\":{\"text\":\"text, watermark\",\"clip\":[\"4\",1]},\"class_type\":\"CLIPTextEncode\",\"_meta\":{\"title\":\"CLIP文本编码\"}},\"8\":{\"inputs\":{\"samples\":[\"3\",0],\"vae\":[\"4\",2]},\"class_type\":\"VAEDecode\",\"_meta\":{\"title\":\"VAE解码\"}},\"9\":{\"inputs\":{\"filename_prefix\":\"ComfyUI\",\"images\":[\"8\",0]},\"class_type\":\"SaveImage\",\"_meta\":{\"title\":\"保存图像\"}},\"10\":{\"inputs\":{\"seed\":-1},\"class_type\":\"Seed (rgthree)\",\"_meta\":{\"title\":\"Seed (rgthree)\"}}}",
  "description": "第一个工作流",
  "path": "workflow/test",
  "workflow": ""
}
```

**注意：** postman或者swagger调用时，请用`\"`代替`"`，避免json解析报错


1. 创建端点参数
  - 查看创建的端点id,swagger `http://127.0.0.1:9518/swagger/index.html#/Endpoint/post_endpoint_index`
  - 创建端点参数 `http://127.0.0.1:9518/swagger/index.html#/EndpointParam/post_endpointParam_set`
  ```HTTP
  POST /endpointParam/set HTTP/1.1
Host: 127.0.0.1:9518
Content-Type: application/json
Content-Length: 627

{
    "endpoint_id": 1,
    "items": [
        {
            "param_name": "提示词",
            "param_key": "prompt",
            "param_type": "string",
            "json_key": "6.inputs.text",
            "description": ""
        },
        {
            "param_name": "宽度",
            "param_key": "width",
            "param_type": "int",
            "json_key": "5.inputs.width",
            "description": ""
        },
        {
            "param_name": "高度",
            "param_key": "height",
            "param_type": "int",
            "json_key": "5.inputs.height",
            "description": ""
        }
    ]
}
  ```
4. 获取http路由，同步调用http接口
    - 获取http路由 `http://127.0.0.1:9518/swagger/index.html#/Endpoint/post_endpoint_index` 同步路由字段`sync_path`，暂不支持异步调用，字段`path`
    - 同步调用http接口，参照步骤2对应创建字段`path`，示例对应api为`http://127.0.0.1:9518/workflow/test/sync`
    - 
    ```HTTP
    POST /workflow/test/sync HTTP/1.1
    Host: 127.0.0.1:9518
    Content-Type: application/json
    Content-Length: 180

    {
        "uid": "1a9828d7-2e83-49f0-9b28-e64443e25a90",
        "prompt": "beautiful scenery nature glass bottle landscape, , purple galaxy bottle,",
        "width": 400,
        "height": 400
    }
    ```
    
    **注意：** uid为comfy应用的uid，必填项，其余参数为步骤3创建的端点参数

5. 获取http路由，异步调用http接口
  - 获取http路由 `http://127.0.0.1:9518/swagger/index.html#/Endpoint/post_endpoint_index` 异步路由字段`path`，字段`path`
  - 异步调用http接口，参照步骤2对应创建字段`path`，示例对应api为`http://127.0.0.1:9518/workflow/test`
  - 
  ```HTTP
  POST /workflow/test HTTP/1.1
  Host: 127.0.0.1:9518
  Content-Type: application/json
  Content-Length: 180

  {
      "uid": "1a9828d7-2e83-49f0-9b28-e64443e25a90",
      "prompt": "beautiful scenery nature glass bottle landscape, , purple galaxy bottle,",
      "width": 400,
      "height": 400,
      "callback_url": "http://127.0.0.1:9999/workflow/test/callback"

  }
  ```
  
  **注意：** uid为comfy应用的uid必填项，callback_url为自定义的POST回调地址、非必填项（优先于endpoint创建时写入的callback_url），其余参数为步骤3创建的端点参数
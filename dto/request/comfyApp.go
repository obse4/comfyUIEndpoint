package request

type ComfyAppCreateRequest struct {
	Addr        string `json:"addr"` // 127.0.0.1:8000
	Description string `json:"description"`
}

type ComfyAppRemoveRequest struct {
	Id int64 `json:"id"`
}

type ComfyAppUpdateRequest struct {
	Id          int64  `json:"id"`
	Addr        string `json:"addr"`
	Description string `json:"description"`
}
type ComfyAppIndexRequest struct {
	Page        int    `json:"page"`
	PageSize    int    `json:"page_size"`
	Addr        string `json:"addr"`
	Description string `json:"description"`
}

type ComfyAppRestartWsRequest struct {
	Id int64 `json:"id"`
}

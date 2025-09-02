package request

type EndpointCreateRequest struct {
	Path        string `json:"path"`        // 路径
	Description string `json:"description"` // 描述
	ApiJson     string `json:"api_json"`    // comfy api版本json
	Workflow    string `json:"workflow"`    // comfy 工作流
}

type EndpointUpdateRequest struct {
	Id          int64  `json:"id"`
	Description string `json:"description"` // 描述
	ApiJson     string `json:"api_json"`    // comfy api版本json
	Workflow    string `json:"workflow"`    // comfy 工作流
}

type EndpointIndexRequest struct {
	Page        int    `json:"page"`
	PageSize    int    `json:"page_size"`
	Path        string `json:"path"`
	Description string `json:"description"`
}

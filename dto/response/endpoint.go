package response

import "comfyui_endpoint/model"

type EndpointIndexResponse struct {
	Data  []model.Endpoint `json:"data"`
	Total int64            `json:"total"`
}

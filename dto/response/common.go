package response

type CommonResponse struct {
	Code    int         `json:"code" example:"200"`
	Message string      `json:"message" example:"ok"`
	Data    interface{} `json:"data"`
}

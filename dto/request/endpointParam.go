package request

type EndpointParamSetRequest struct {
	EndpointId int64                      `json:"endpoint_id"`
	Items      []EndpointParamItemRequest `json:"items"`
}

type EndpointParamItemRequest struct {
	ParamName   string `json:"param_name"`
	ParamKey    string `json:"param_key"`
	Description string `json:"description"`
	ParamType   string `json:"param_type"` // string/int/float
	JsonKey     string `json:"json_key"`
}

type EndpointParamFindRequest struct {
	EndpointId int64 `json:"endpoint_id"`
}

package model

import "time"

type Endpoint struct {
	Id          int64     `json:"id" gorm:"column:id;primary_key;NOT NULL;comment:'编号'"`
	Path        string    `json:"path" gorm:"column:path;NOT NULL;comment:'路径'"`
	SyncPath    string    `json:"sync_path" gorm:"column:sync_path;comment:'同步路径'"`
	Description string    `json:"description" gorm:"column:description;comment:'描述'"`
	ApiJson     string    `json:"api_json" gorm:"column:api_json;type:longtext;comment:'api_json'"` // comfyUI api用工作流
	Workflow    string    `json:"workflow" gorm:"column:workflow;type:longtext;comment:'工作流'"`
	CallbackUrl string    `json:"callback_url" gorm:"column:callback_url;comment:'回调地址'"` // http://127.0.0.1:8080/api/v1/callback
	CreatedAt   time.Time `json:"created_at" gorm:"column:created_at;"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"column:updated_at;"`
}

func (Endpoint) TableName() string {
	return "endpoints"
}

type EndpointParam struct {
	Id          int64     `json:"id" gorm:"column:id;primary_key;NOT NULL;comment:'编号'"`
	EndpointId  int64     `json:"endpoint_id" gorm:"column:endpoint_id;NOT NULL;comment:'endpoint_id'"`
	ParamName   string    `json:"param_name" gorm:"column:param_name;comment:'参数名'"`
	ParamKey    string    `json:"param_key" gorm:"column:param_key;NOT NULL;comment:'参数key'"`
	Description string    `json:"description" gorm:"column:description;comment:'描述'"`
	ParamType   string    `json:"param_type" gorm:"column:param_type;NOT NULL;comment:'参数类型'"` // string, int, float
	JsonKey     string    `json:"json_key" gorm:"column:json_key;NOT NULL;comment:'json key'"` // 用于json解析的key
	CreatedAt   time.Time `json:"created_at" gorm:"column:created_at;"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"column:updated_at;"`
}

func (EndpointParam) TableName() string {
	return "endpoint_params"
}

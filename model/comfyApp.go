package model

import "time"

type ComfyApp struct {
	Id          int64     `json:"id" gorm:"column:id;primary_key;NOT NULL;comment:'编号'"`
	Uid         string    `json:"uid" gorm:"column:uid;NOT NULL;comment:'uuid'"`
	Addr        string    `json:"addr" gorm:"column:addr;NOT NULL;comment:'地址'"`
	Description string    `json:"description" gorm:"column:description;comment:'描述'"`
	CreatedAt   time.Time `json:"created_at" gorm:"column:created_at;"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"column:updated_at;"`
}

func (ComfyApp) TableName() string {
	return "comfy_apps"
}

type ComfyAppInfo struct {
	Id         int64     `json:"id" gorm:"column:id;primary_key;NOT NULL;comment:'编号'"`
	Uid        string    `json:"uid" gorm:"column:uid;NOT NULL;comment:'uuid'"`
	WsStatus   string    `json:"ws_status" gorm:"column:ws_status;default:'close';comment:'websocket状态'"`
	RunningNum int64     `json:"running_num" gorm:"column:running_num;NOT NULL;comment:'运行中数量'"`
	PendingNum int64     `json:"pending_num" gorm:"column:pending_num;NOT NULL;comment:'等待中数量'"`
	CreatedAt  time.Time `json:"created_at" gorm:"column:created_at;"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"column:updated_at;"`
}

func (ComfyAppInfo) TableName() string {
	return "comfy_app_infos"
}

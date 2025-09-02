package service

import (
	"comfyui_endpoint/config"

	"gorm.io/gorm"
)

func SqliteDb() *gorm.DB {
	return config.Global.Sqlite.Conn
}

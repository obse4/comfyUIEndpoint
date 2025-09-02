package database

import (
	"comfyui_endpoint/model"

	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		model.ComfyApp{},
		model.ComfyAppInfo{},
		model.Endpoint{},
		model.EndpointParam{},
	)
}

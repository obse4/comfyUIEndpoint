package database

import (
	"comfyui_endpoint/config"
	"comfyui_endpoint/logger"
	"comfyui_endpoint/model"

	"github.com/glebarez/sqlite" // 使用 glebarez/sqlite 驱动
	"gorm.io/gorm"
	log "gorm.io/gorm/logger"
)

func InitSqlite(config *config.SqliteConfig) error {
	var err error

	config.Conn, err = gorm.Open(sqlite.Open(config.Db), &gorm.Config{
		Logger: log.Default.LogMode(log.Info),
	})

	if err != nil {
		return err
	}

	// 获取底层的 sqlDB
	sqlDB, err := config.Conn.DB()
	if err != nil {
		return err
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(10)  // 设置空闲连接池中的最大连接数
	sqlDB.SetMaxOpenConns(100) // 设置打开数据库连接的最大数量

	// 自动迁移数据库结构
	config.Conn.AutoMigrate(
		&model.ComfyApp{},
		&model.ComfyAppInfo{},
		&model.Endpoint{},
		&model.EndpointParam{},
	)

	logger.Debug("InitSqlite Success")
	return nil
}

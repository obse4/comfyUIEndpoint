package config

import (
	"embed"

	"gopkg.in/yaml.v3"
)

var (
	Global Config
	//go:embed local.config.yaml
	configFs embed.FS
)

func InitConfig() {
	// 读取嵌入的配置文件
	data, err := configFs.ReadFile("local.config.yaml")

	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(data, &Global)

	if err != nil {
		panic(err)
	}
}

package config

import (
	"github.com/gomodule/redigo/redis"
	"gorm.io/gorm"
)

type Config struct {
	HttpServer HttpServerConfig `yaml:"httpServer"`
	Log        LogConfig        `yaml:"log"`
	Sqlite     SqliteConfig     `yaml:"sqlite"`
}

type SqliteConfig struct {
	Conn *gorm.DB
	Db   string `yaml:"db"`
}

type HttpServerConfig struct {
	Mode string `yaml:"mode"`
}

type LogConfig struct {
	Level int `yaml:"level"`
}

type RedisConfig struct {
	Name     string      `yaml:"name"`     // 自定义名称
	Url      string      `yaml:"url"`      // url连接
	Port     string      `yaml:"port"`     // 端口
	Password string      `yaml:"password"` // 密码 非必填
	Db       int         `yaml:"db"`
	Pool     *redis.Pool // redis连接池
}

type MysqlConfig struct {
	Username string   `yaml:"username"` // 用户名
	Password string   `yaml:"password"` // 密码
	Database string   `yaml:"database"` // 数据库
	Url      string   `yaml:"url"`      // url地址
	Port     string   `yaml:"port"`     // 端口
	Db       *gorm.DB // 数据库指针
}

package tool

import (
	"os"
	"bufio"
	"encoding/json"
)

type Config struct {
	AppName string `json:"appName"`
	AppMode string `json:"appMode"`
	AppHost string `json:"appHost"`
	AppPort string `json:"appPort"`
	Sms     SmsConfig `json:"sms"`
	Database DatabaseConfig `json:"database"`
	RedisConfig RedisConfig    `json:"redis_config"`
}

type SmsConfig struct {
	SignName     string `json:"sign_name"`
	TemplateCode string `json:"template_code"`
	RegionId     string `json:"region_id"`
	AppKey       string `json:"app_key"`
	AppSecret    string `json:"app_secret"`
}

type DatabaseConfig struct {
	Driver   string `json:"driver"`
	User     string `json:"user"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	DbName   string `json:"db_name"`
	Charset  string `json:"charset"`
	ShowSql  bool   `json:"show_sql"`
}

//Redis属性定义
type RedisConfig struct {
	Addr     string `json:"addr"`
	Port     string `json:"port"`
	Password string `json:"password"`
	Db       int    `json:"db"`
}

var _cfg * Config = nil

func GetConfig() *Config {
	return _cfg
}

// 读取配置文件并解析返回
func ParseConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		panic(err)
	}
	reader := bufio.NewReader(file)
	decoder := json.NewDecoder(reader)
	if err = decoder.Decode(&_cfg); err != nil{	
		return nil, err
	}
	return _cfg, nil
}
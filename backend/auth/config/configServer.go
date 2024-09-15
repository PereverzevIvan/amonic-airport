package config

import (
	"encoding/json"
	"flag"
	"os"
	"time"
)

type ConfigServer struct {
	Port int `json:"server_port"`
}

type ConfigDatabase struct {
	DBType     string `json:"db_type"`
	Host       string `json:"host"`
	Port       int    `json:"db_port"`
	DBName     string `json:"db_name"`
	AdminName  string `json:"admin_name"`
	DBPassword string `json:"db_password"`
	SSL        string `json:"ssl"`
}

type ConfigJWT struct {
	SecretKey                 string        `json:"secret_key"`
	AccessTokenExpirationStr  string        `json:"access_token_expiration"`
	RefreshTokenExpirationStr string        `json:"refresh_token_expiration"`
	AccessTokenExpiration     time.Duration `json:"-"`
	RefreshTokenExpiration    time.Duration `json:"-"`
}

type Config struct {
	ConfigServer   `json:"server"`
	ConfigDatabase `json:"database"`
	ConfigJWT      `json:"jwt"`
}

func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "./config/config_dev.json", "path to config file") // Сначала пробуем считать путь из консоли
	flag.Parse()

	return res
}

func MustLoadConfig() Config {
	path := fetchConfigPath()
	if path == "" {
		panic("Путь до конфига пустой")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("Файл не найден")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		panic("Не удалось прочитать файл: " + path)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		panic("Не удалось прочитать json файл: " + err.Error())
	}

	cfg.AccessTokenExpiration, err = time.ParseDuration(cfg.AccessTokenExpirationStr)
	if err != nil {
		panic("Не удалось прочитать AccessTokenExpiration: " + err.Error())
	}

	cfg.RefreshTokenExpiration, err = time.ParseDuration(cfg.RefreshTokenExpirationStr)
	if err != nil {
		panic("Не удалось прочитать RefreshTokenExpiration: " + err.Error())
	}

	return cfg
}

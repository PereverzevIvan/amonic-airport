package config

import (
	"encoding/json"
	"flag"
	"os"
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

type Config struct {
	ConfigServer   `json:"server"`
	ConfigDatabase `json:"database"`
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
		panic("Не удалось прочитать файл")
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		panic("Не удалось прочитать json файл")
	}

	return cfg
}

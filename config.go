package main

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	Databases []DatabaseConfig `json:"databases"`
	Qiniu     QiniuConfig      `json:"qiniu"`
	Cron      string           `json:"cron"`
}

func loadConfig(file *string) Config {
	fileRaw, err := os.ReadFile(*file)
	if err != nil {
		log.Fatalf("Read config file failed: %s", err)
	}

	var config Config
	if e := json.Unmarshal(fileRaw, &config); e != nil {
		log.Fatalf("Parse config file failed: %s", e)
	}
	return config
}

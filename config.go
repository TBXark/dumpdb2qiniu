package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type Config struct {
	Databases []DatabaseConfig `json:"databases"`
	Qiniu     QiniuConfig      `json:"qiniu"`
	Cron      string           `json:"cron"`
}

func loadConfig(path string) (*Config, error) {
	var body []byte
	var err error

	if strings.HasPrefix(path, "http") {
		resp, httpErr := http.Get(path)
		if httpErr != nil {
			return nil, fmt.Errorf("failed to fetch config: %w", httpErr)
		}
		defer resp.Body.Close()
		body, err = io.ReadAll(resp.Body)
	} else {
		body, err = os.ReadFile(path)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}
	var conf Config
	err = json.Unmarshal(body, &conf)
	if err != nil {
		return nil, err
	}
	return &conf, nil
}

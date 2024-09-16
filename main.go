package main

import (
	"flag"
	"fmt"
	"github.com/robfig/cron/v3"
	"log"
	"path/filepath"
	"strings"
)

type Uploader interface {
	upload(key, filePath string) error
}

type Dumper interface {
	Dump() (string, error)
}

func start(conf *Config) {
	uploader := NewQiniuUploader(&conf.Qiniu)
	for _, db := range conf.Databases {
		dump := NewMysqlDump(&db)
		file, err := dump.Dump()
		if err != nil {
			log.Printf("Backup database %s failed: %s", db.Database, err)
			continue
		}
		fileName := filepath.Base(file)
		key := fmt.Sprintf("%s/%s/%s/%s", strings.TrimRight(conf.Qiniu.FilePrefix, "/"), db.Host, db.User, fileName)
		err = uploader.upload(key, file)
		if err != nil {
			log.Printf("Upload to Qiniu failed: %s", err)
			continue
		}
	}
}

func main() {
	conf := flag.String("config", "config.json", "config file path")
	flag.Parse()
	config := loadConfig(conf)

	if config.Cron != "" {
		task := cron.New()
		_, e := task.AddFunc(config.Cron, func() {
			start(&config)
		})
		if e != nil {
			log.Fatalf("Run cron task failed: %s", e)
		}
		task.Run()
	} else {
		start(&config)
	}
}

package main

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/robfig/cron/v3"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

type DatabaseConfig struct {
	Host     string `json:"host"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
}

type QiniuConfig struct {
	AccessKey  string `json:"access_key"`
	SecretKey  string `json:"secret_key"`
	Bucket     string `json:"bucket"`
	FilePrefix string `json:"file_prefix"`
}

type Config struct {
	Databases []DatabaseConfig `json:"databases"`
	Qiniu     QiniuConfig      `json:"qiniu"`
	Cron      string           `json:"cron"`
}

func backup(db *DatabaseConfig, qiniu *QiniuConfig) {
	fmt.Printf("Start backup database %s\n", db.Database)
	dumpCmd := fmt.Sprintf("mysqldump -h %s -u %s -p%s %s", db.Host, db.User, db.Password, db.Database)
	cmd := exec.Command("bash", "-c", dumpCmd)
	var out bytes.Buffer
	cmd.Stdout = &out
	if e := cmd.Run(); e != nil {
		log.Fatalf("Dump database %s failed: %s", db.Database, e)
	}

	// 压缩导出的数据
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	if _, e := gz.Write(out.Bytes()); e != nil {
		log.Fatalf("Compress database %s failed: %s", db.Database, e)
	}
	if e := gz.Close(); e != nil {
		log.Fatalf("Compress database %s failed: %s", db.Database, e)
	}

	fileName := fmt.Sprintf("%s_%s.gz", db.Database, time.Now().Format("20060102150405"))
	tmpFile, err := os.CreateTemp("", fileName)
	if err != nil {
		log.Fatalf("Create temp file failed: %s", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, e := tmpFile.Write(b.Bytes()); e != nil {
		log.Fatalf("Write temp file failed: %s", e)
	}
	if e := tmpFile.Close(); e != nil {
		log.Fatalf("Close temp file failed: %s", e)
	}

	// 上传到七牛云
	putPolicy := storage.PutPolicy{
		Scope: qiniu.Bucket,
	}
	mac := qbox.NewMac(qiniu.AccessKey, qiniu.SecretKey)
	upToken := putPolicy.UploadToken(mac)

	cfg := storage.Config{}
	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}

	key := fmt.Sprintf("%s/%s/%s/%s", strings.TrimRight(qiniu.FilePrefix, "/"), db.Host, db.User, fileName)
	err = formUploader.PutFile(context.Background(), &ret, upToken, key, tmpFile.Name(), nil)
	if err != nil {
		log.Fatalf("Upload to qiniu failed: %s", err)
	}
	fmt.Printf("Backup database %s@%s %s success\n", db.User, db.Host, db.Database)
}

func start(cfg *Config) {
	for _, db := range cfg.Databases {
		backup(&db, &cfg.Qiniu)
	}
}

func main() {
	cfgParam := flag.String("config", "config.json", "config file path")
	flag.Parse()
	fileRaw, err := os.ReadFile(*cfgParam)
	if err != nil {
		log.Fatalf("Read config file failed: %s", err)
	}
	var cfg Config
	if e := json.Unmarshal(fileRaw, &cfg); e != nil {
		log.Fatalf("Parse config file failed: %s", e)
	}
	if cfg.Cron != "" {
		task := cron.New()
		_, e := task.AddFunc(cfg.Cron, func() {
			start(&cfg)
		})
		if e != nil {
			log.Fatalf("Add cron task failed: %s", e)
		}
		task.Run()
	} else {
		start(&cfg)
	}

}

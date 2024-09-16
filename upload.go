package main

import (
	"context"
	"fmt"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
)

type QiniuConfig struct {
	AccessKey  string `json:"access_key"`
	SecretKey  string `json:"secret_key"`
	Bucket     string `json:"bucket"`
	FilePrefix string `json:"file_prefix"`
}

type QiniuUploader struct {
	conf *QiniuConfig
}

func NewQiniuUploader(conf *QiniuConfig) Uploader {
	return &QiniuUploader{conf: conf}
}

func (q *QiniuUploader) upload(key, filePath string) error {

	putPolicy := storage.PutPolicy{
		Scope: q.conf.Bucket,
	}
	mac := qbox.NewMac(q.conf.AccessKey, q.conf.SecretKey)
	upToken := putPolicy.UploadToken(mac)

	cfg := storage.Config{}
	ret := storage.PutRet{}
	formUploader := storage.NewFormUploader(&cfg)

	err := formUploader.PutFile(context.Background(), &ret, upToken, key, filePath, nil)
	if err != nil {
		return fmt.Errorf("upload to qiniu failed: %w", err)
	}

	return nil
}

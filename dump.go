package main

import (
	"compress/gzip"
	"fmt"
	"os"
	"os/exec"
	"time"
)

type DatabaseConfig struct {
	Host     string `json:"host"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
}

type MysqlDump struct {
	conf *DatabaseConfig
}

func NewMysqlDump(conf *DatabaseConfig) Dumper {
	return &MysqlDump{conf: conf}
}

func (m *MysqlDump) Dump() (string, error) {
	db := m.conf
	fmt.Printf("Start backup database %s\n", db.Database)

	dumpCmd := fmt.Sprintf("mysqldump -h %s -u %s -p%s %s", db.Host, db.User, db.Password, db.Database)
	cmd := exec.Command("bash", "-c", dumpCmd)

	fileName := fmt.Sprintf("%s_%s.gz", db.Database, time.Now().Format("20060102150405"))
	tmpFile, err := os.CreateTemp("", fileName)

	if err != nil {
		return "", fmt.Errorf("create temp file failed: %w", err)
	}
	defer tmpFile.Close()

	gzWriter := gzip.NewWriter(tmpFile)
	defer gzWriter.Close()
	cmd.Stdout = gzWriter

	if e := cmd.Run(); e != nil {
		_ = os.Remove(tmpFile.Name())
		return "", fmt.Errorf("dump database %s failed: %w", db.Database, e)
	}

	if e := gzWriter.Close(); e != nil {
		_ = os.Remove(tmpFile.Name())
		return "", fmt.Errorf("close gzip writer failed: %w", e)
	}

	if e := tmpFile.Close(); e != nil {
		_ = os.Remove(tmpFile.Name())
		return "", fmt.Errorf("close temp file failed: %w", e)
	}

	return tmpFile.Name(), nil
}

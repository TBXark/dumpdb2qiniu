# DumpDB2Qiniu

这是一个自动备份数据库到七牛云的工具，当前支持所有能使用mysqldump备份的数据库。

### 使用
```shell
dumpdb2qiniu -c config.json
```

### 安装
1. 手动编译
```shell
go install github.com/TBXark/dumpdb2qiniu
```

### 配置文件

```json
{
  "databases": [
    {
      "host": "localhost",
      "user": "root",
      "password": "....",
      "database": "dbname"
    }
  ], // 需要备份的数据库，支持多个
  "qiniu": {
    "access_key": "-e",
    "secret_key": "",
    "bucket": "sqldump", // 为了保证数据安全，请使用私有空间
    "file_prefix": "backup/serverA/"
  },
  "cron": "" // 需要定时执行的时间，格式为cron表达式
}
```
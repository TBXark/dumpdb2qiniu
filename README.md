# dumpdb2qiniu

**dumpdb2qiniu** is a tool that automatically backs up databases to Qiniu Cloud. It currently supports backing up all databases that can be backed up using mysqldump.

## Installation

### Manual Build
```bash
go install github.com/TBXark/dumpdb2qiniu@latest
```

### Docker
```bash
docker run -d --name dumpdb2qiniu -v /path/to/config.json:/config.json ghcr.io/tbxark/dumpdb2qiniu:latest
```

## Usage

```bash
dumpdb2qiniu -config /path/to/config.json
```

## Configuration

```json5
{
  "databases": [
    {
      "host": "localhost",
      "user": "root",
      "password": "....",
      "database": "dbname"
    }
  ],  // Multiple databases are supported
  "qiniu": {
    "access_key": "-e",
    "secret_key": "",
    "bucket": "sqldump", // For security, please use a private bucket
    "file_prefix": "backup/serverA/"
  },
  "cron": "" // Cron expression for scheduled execution
}
```

## License

**dumpdb2qiniu**  is licensed under the MIT License. See the [LICENSE](./LICENSE) file for more details.
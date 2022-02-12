# SOGNO file-service
This service stores static grid data, like CIM CGMES network descriptions, for other SOGNO services.

## Development

### Compiling

```bash
$ go mod tidy
$ go build
```

### Generating OpenAPI docs

```bash
# Ensure your Go bin directory is on your path (default: ~/go/bin)
swag init
```

### Configuring

```bash
# For example:
mkdir -p ~/.config/sogno-file-service
echo '{"minio_endpoint": "s3.amazonaws.com", "minio_bucket": "'$SOGNO_FILE_SERVICE_BUCKET'"}' > ~/.config/sogno-file-service/config.json
```

### Running

```bash
go run main.go
```

# SOGNO file-service
This service stores static grid data, like CIM CGMES network descriptions, for other SOGNO services.

## Development

### Compiling

```bash
$ go mod tidy
$ go build
```

### Docker build / run

The docker build uses the file "minio.config" to store the s3 details.
You will need to make sure this matches your setup.

```bash
$ docker build -t sogno-file-service .
$ docker run -p 8080:8080 sogno-file-service
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

### Documentation

Visit localhost:8080 in your web browser to view the HTML version of
the API documentation. You can also view `docs/swagger.yaml` in the
repo.

### Testing

Currently, the only tests are integration tests and they require a running
S3-compatible object storage server.

```bash
go test
```

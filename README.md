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

### Running

```bash
go run main.go
```

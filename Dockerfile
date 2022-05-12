
FROM golang:latest

RUN mkdir /usr/src/app
WORKDIR /usr/src/app
COPY main.go go.mod /usr/src/app/
COPY routes /usr/src/app/routes/
COPY config /usr/src/app/config/
COPY docs /usr/src/app/docs/
COPY api /usr/src/app/api/
COPY file /usr/src/app/file/
RUN mkdir -p /usr/src/app/.config/sogno-file-service
COPY minio.config /usr/src/app/.config/sogno-file-service/config.json
RUN go mod tidy
RUN go build
RUN useradd app -d /usr/src/app
RUN chown app -R /usr/src/app
USER app

EXPOSE 8080

CMD go run main.go


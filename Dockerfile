FROM golang:1.20 as builder
WORKDIR /source
COPY . /source
RUN go build -o app .

FROM ubuntu:noble
COPY --from=builder /source/app /app
RUN chmod +x /app

ENV APP_NAME="app"
ENV APP_VERSION="v0.0.1"
ENV LOG_LEVEL="INFO"
ENV ADMIN_NAME="admin"
ENV REST_LISTEN_ADDR="0.0.0.0:8090"
ENV WEBSOCKET_LISTEN_ADDR="0.0.0.0:8091"
ENV DSN="postgres://postgres:password@localhost:5432/thinsectiondev?sslmode=disable"
ENV BYPASS_AUTH=true

CMD ["/app"]
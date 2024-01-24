FROM golang:1.20 as builder
WORKDIR /source
COPY . /source
RUN go build -o app .

FROM ubuntu:noble
COPY --from=builder /source/app /app
RUN chmod +x /app
CMD ["/app"]
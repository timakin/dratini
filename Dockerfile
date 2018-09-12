FROM golang:1.11.0 AS builder
ADD . ~/github.com/timakin/dratini
WORKDIR ~/github.com/timakin/dratini
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app

FROM alpine:latest as runner
COPY --from=builder ~/github.com/timakin/dratini/app /usr/local/bin/app


ARG SERVICE=api-gateway
ARG TARGET_DIR=/app
ARG GOBIN=/.gobin

FROM docker.io/library/golang:1.22.3-alpine as builder
ARG TARGET_DIR
ARG GOBIN
ARG ISLOCAL
ARG SERVICE
ENV SERVICE=$SERVICE

# gcc and musl-dev is required to support golangci-lint (install-tools command) 
RUN apk add --update make git musl-dev gcc linux-headers

ENV GO111MODULE=on
WORKDIR /build

COPY go.mod go.mod
COPY go.sum go.sum
COPY configs/initial.yml $TARGET_DIR/config.yml

RUN go mod download
COPY . .
RUN go build -o ${TARGET_DIR}/${SERVICE} cmd/app/main.go

# Main image
FROM alpine:3.14
ARG TARGET_DIR
ARG GOBIN
ARG SERVICE
ENV SERVICE=$SERVICE

RUN adduser -D app && mkdir -p /app && chown -R app /app
USER app

WORKDIR /app/
COPY --from=builder ${TARGET_DIR}/${SERVICE} ${SERVICE}
COPY --from=builder ${TARGET_DIR}/config.yml ./configs/config.yml

EXPOSE 8080
ENTRYPOINT /app/${SERVICE}

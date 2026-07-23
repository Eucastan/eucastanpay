FROM golang:1.25 AS builder

WORKDIR /src

# Dependency files
COPY services/api-gateway/go.mod ./services/api-gateway/go.mod
COPY services/api-gateway/go.sum ./services/api-gateway/go.sum

# Common module
COPY common ./common
WORKDIR /src/services/api-gateway

RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# Account source
COPY services/api-gateway ./

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-w -s" -o /api-gateway ./cmd/main.go

FROM alpine:3.22

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /api-gateway .

RUN adduser -D -u 1001 appuser

USER appuser

ENTRYPOINT ["./api-gateway"]
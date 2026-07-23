FROM golang:1.25 AS builder

WORKDIR /src

# Dependency files
COPY services/admin/go.mod ./services/admin/go.mod
COPY services/admin/go.sum ./services/admin/go.sum

# Common module
COPY common ./common
WORKDIR /src/services/admin

RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# admin source
COPY services/admin ./

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-w -s" -o /admin ./cmd/admin

FROM alpine:3.22

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /admin .

RUN adduser -D -u 1001 appuser

USER appuser

ENTRYPOINT ["./admin"]
FROM golang:1.25 AS builder

WORKDIR /src

# Dependency files
COPY services/notification/go.mod ./services/notification/go.mod
COPY services/notification/go.sum ./services/notification/go.sum

# Common module
COPY common ./common
WORKDIR /src/services/notification

RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# notification source
COPY services/notification ./

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-w -s" -o /notification ./cmd/notification

FROM alpine:3.22

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /notification .

RUN adduser -D -u 1001 appuser

USER appuser

ENTRYPOINT ["./notification"]
FROM golang:1.25 AS builder

WORKDIR /src

# Dependency files
COPY services/account/go.mod ./services/account/go.mod
COPY services/account/go.sum ./services/account/go.sum

# Common module
COPY common ./common
WORKDIR /src/services/account

RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# Account source
COPY services/account ./

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-w -s" -o /account ./cmd/account

FROM alpine:3.22

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /account .

RUN adduser -D -u 1001 appuser

USER appuser

ENTRYPOINT ["./account"]
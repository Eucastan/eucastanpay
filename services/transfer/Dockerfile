FROM golang:1.25 AS builder

WORKDIR /src

# Dependency files
COPY services/transfer/go.mod ./services/transfer/go.mod
COPY services/transfer/go.sum ./services/transfer/go.sum

# Common module
COPY common ./common
WORKDIR /src/services/transfer

RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# transfer source
COPY services/transfer ./

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-w -s" -o /transfer ./cmd/transfer

FROM alpine:3.22

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /transfer .

RUN adduser -D -u 1001 appuser

USER appuser

ENTRYPOINT ["./transfer"]
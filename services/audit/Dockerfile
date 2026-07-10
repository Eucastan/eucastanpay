FROM golang:1.25 AS builder

WORKDIR /src

# Dependency files
COPY services/audit/go.mod ./services/audit/go.mod
COPY services/audit/go.sum ./services/audit/go.sum

# Common module
COPY common ./common
WORKDIR /src/services/audit

RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# audit source
COPY services/audit ./

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-w -s" -o /audit ./cmd/audit

FROM alpine:3.22

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /audit .

RUN adduser -D -u 1001 appuser

USER appuser

ENTRYPOINT ["./audit"]
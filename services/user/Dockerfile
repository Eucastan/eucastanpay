FROM golang:1.25 AS builder

WORKDIR /src

# Dependency files
COPY services/user/go.mod ./services/user/go.mod
COPY services/user/go.sum ./services/user/go.sum

# Common module
COPY common ./common
WORKDIR /src/services/user

RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# user source
COPY services/user ./

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-w -s" -o /user ./cmd/user

FROM alpine:3.22

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /user .

RUN adduser -D -u 1001 appuser

USER appuser

ENTRYPOINT ["./user"]
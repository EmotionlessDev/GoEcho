# syntax=docker/dockerfile:1

ARG GO_VERSION=1.27.1
ARG ALPINE_VERSION=3.22

FROM golang:${GO_VERSION}-alpine AS build

WORKDIR /src

# Зависимости отдельным слоем: пересобираются только при изменении go.mod/go.sum,
# а не на каждую правку кода.
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .

ARG SERVICE

# CGO выключен — получаем статический бинарь, который едет в alpine без зависимостей.
# Кеш сборки смонтирован, поэтому повторная сборка соседнего сервиса почти мгновенна.
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/service ./cmd/${SERVICE}

FROM alpine:${ALPINE_VERSION} AS runtime

RUN apk add --no-cache ca-certificates tzdata wget \
 && adduser -D -u 10001 app \
 && mkdir -p /data && chown app:app /data

COPY --from=build /out/service /usr/local/bin/service

USER app

ENTRYPOINT ["/usr/local/bin/service"]

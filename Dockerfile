ARG APP_IMAGE=golang:1.17
FROM ${APP_IMAGE} as build

ARG BUILD_VERSION=docker
ARG TARGETOS
ARG TARGETARCH
ARG GO_PATH="/src/.go"

ENV CGO_ENABLED 0
ENV GOOS=$TARGETOS
ENV GOARCH=$TARGETARCH
ENV GOPATH=$GO_PATH
ENV GOPROXY https://proxy.golang.org,direct
ENV GOSUMDB off
ENV TZ Europe/Moscow

WORKDIR /src
COPY . .
RUN make build BUILD_VERSION=$BUILD_VERSION

FROM scratch
ENV TZ Europe/Moscow
WORKDIR /app
COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /src/bin bin
COPY configs configs
ENTRYPOINT ["bin/app"]

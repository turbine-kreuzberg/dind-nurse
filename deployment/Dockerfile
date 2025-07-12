# dependencies download
FROM docker:28.3.2@sha256:684bc226a69eb91582738f196e59e234355420ee592abea7685e4127afbce645 AS dependencies
ENV BUILDX_VERSION=v0.25.0
RUN set -eux; \
    wget -q -O /docker-buildx "https://github.com/docker/buildx/releases/download/${BUILDX_VERSION:?}/buildx-${BUILDX_VERSION:?}.linux-amd64"; \
    chmod a+x /docker-buildx

# build environment ###########################################
FROM golang:1.24.4-alpine@sha256:68932fa6d4d4059845c8f40ad7e654e626f3ebd3706eef7846f319293ab5cb7a AS build-env
WORKDIR /app

RUN apk add --no-cache ca-certificates docker

# buildx
COPY --from=dependencies /docker-buildx /root/.docker/cli-plugins/docker-buildx

# entrypoint
RUN apk add --no-cache ca-certificates docker entr build-base
COPY ./deployment/entrypoint.sh /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]

# dependencies
COPY go.mod .
COPY go.sum .
RUN go mod download

# server
COPY main.go .
COPY pkg ./pkg

RUN go install .

# debug environment ###########################################
FROM build-env AS debug-env
RUN go install github.com/go-delve/delve/cmd/dlv@v1.8.0
COPY ./deployment/entrypoint-delve.sh /entrypoint.sh
RUN go install -gcflags "all=-N -l" .

# production image ############################################
FROM alpine:3.22.0@sha256:8a1f59ffb675680d47db6337b49d22281a139e9d709335b492be023728e11715 AS prod

RUN apk add --no-cache ca-certificates docker

# buildx
COPY --from=dependencies /docker-buildx /root/.docker/cli-plugins/docker-buildx

# executable
COPY --from=build-env /go/bin/dind-nurse /bin/dind-nurse
ENTRYPOINT ["dind-nurse"]

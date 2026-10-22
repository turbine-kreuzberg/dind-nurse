# dependencies download
FROM docker:29.9.1@sha256:3e5756e2b6dad6a5f1f5deacf0a19f9eaa2581b76ae092b2d4a118162e1d4d4e AS dependencies
ENV BUILDX_VERSION=v0.38.1
RUN set -eux; \
    wget -q -O /docker-buildx "https://github.com/docker/buildx/releases/download/${BUILDX_VERSION:?}/buildx-${BUILDX_VERSION:?}.linux-amd64"; \
    chmod a+x /docker-buildx

# build environment ###########################################
FROM golang:1.27.3-alpine@sha256:d530cf60723abf1ffae0cd16bb546b379ec0550d2f3815e5b16ed69dbc00fba4 AS build-env
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
FROM alpine:3.24.3@sha256:b3413572228e34bd1c9245f5565e6cd127005655e8a7899aef33d8e83255c9e6 AS prod

RUN apk add --no-cache ca-certificates docker

# buildx
COPY --from=dependencies /docker-buildx /root/.docker/cli-plugins/docker-buildx

# executable
COPY --from=build-env /go/bin/dind-nurse /bin/dind-nurse
ENTRYPOINT ["dind-nurse"]

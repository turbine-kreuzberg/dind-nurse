# dependencies download
FROM docker:26.1.3@sha256:1ff6086d17277f3bc4aa1da23426c84f39d256fe8b3d54fa54ebb1c05863f58d AS dependencies
ENV BUILDX_VERSION=v0.14.1
RUN set -eux; \
    wget -q -O /docker-buildx "https://github.com/docker/buildx/releases/download/${BUILDX_VERSION:?}/buildx-${BUILDX_VERSION:?}.linux-amd64"; \
    chmod a+x /docker-buildx

# build environment ###########################################
FROM golang:1.22.3-alpine@sha256:7e788330fa9ae95c68784153b7fd5d5076c79af47651e992a3cdeceeb5dd1df0 AS build-env
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
FROM alpine:3.20.0@sha256:77726ef6b57ddf65bb551896826ec38bc3e53f75cdde31354fbffb4f25238ebd AS prod

RUN apk add --no-cache ca-certificates docker

# buildx
COPY --from=dependencies /docker-buildx /root/.docker/cli-plugins/docker-buildx

# executable
COPY --from=build-env /go/bin/dind-nurse /bin/dind-nurse
ENTRYPOINT ["dind-nurse"]

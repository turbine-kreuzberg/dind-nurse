# dependencies download
FROM docker:20.10.20@sha256:b39aa91b8a31462d6e86da84762891f6d2ac5657d34d9bf86fd1d0bd7b60e0d8 AS dependencies
RUN apk add curl
ENV BUILDX_VERSION=v0.8.2
RUN curl --fail -L -o /docker-buildx \
    https://github.com/docker/buildx/releases/download/${BUILDX_VERSION}/buildx-${BUILDX_VERSION}.linux-amd64
RUN chmod a+x /docker-buildx

# build environment ###########################################
FROM golang:1.19.2-alpine@sha256:f3e683657ddf73726b5717c2ff80cdcd9e9efb7d81f77e4948fada9a10dc7257 AS build-env
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
FROM alpine:3.16.2@sha256:bc41182d7ef5ffc53a40b044e725193bc10142a1243f395ee852a8d9730fc2ad AS prod

RUN apk add --no-cache ca-certificates docker

# buildx
COPY --from=dependencies /docker-buildx /root/.docker/cli-plugins/docker-buildx

# executable
COPY --from=build-env /go/bin/dind-nurse /bin/dind-nurse
ENTRYPOINT ["dind-nurse"]

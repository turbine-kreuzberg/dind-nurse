# dependencies download
FROM docker:20.10.17@sha256:386dd216988a21c83233d32b3926300909a111fa73266d035e992280aa3a9e42 AS dependencies
RUN apk add curl
ENV BUILDX_VERSION=v0.8.2
RUN curl --fail -L -o /docker-buildx \
    https://github.com/docker/buildx/releases/download/${BUILDX_VERSION}/buildx-${BUILDX_VERSION}.linux-amd64
RUN chmod a+x /docker-buildx

# build environment ###########################################
FROM golang:1.18.4-alpine@sha256:46f1fa18ca1ec228f7ea4978ad717f0a8c5e51436e7b8efaf64011f7729886df AS build-env
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
FROM alpine:3.16.0@sha256:686d8c9dfa6f3ccfc8230bc3178d23f84eeaf7e457f36f271ab1acc53015037c AS prod

RUN apk add --no-cache ca-certificates docker

# buildx
COPY --from=dependencies /docker-buildx /root/.docker/cli-plugins/docker-buildx

# executable
COPY --from=build-env /go/bin/dind-nurse /bin/dind-nurse
ENTRYPOINT ["dind-nurse"]

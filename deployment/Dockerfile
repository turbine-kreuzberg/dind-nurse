# build environment ###########################################
FROM golang:1.17.6-alpine@sha256:519c827ec22e5cf7417c9ff063ec840a446cdd30681700a16cf42eb43823e27c AS build-env

WORKDIR /app

# entrypoint
RUN apk add --no-cache entr docker build-base
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
FROM alpine:3.15.0@sha256:21a3deaa0d32a8057914f36584b5288d2e5ecc984380bc0118285c70fa8c9300 AS prod
RUN apk add --no-cache ca-certificates

# executable
COPY --from=build-env /go/bin/dind-nurse /bin/dind-nurse
ENTRYPOINT ["dind-nurse"]

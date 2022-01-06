# build environment ###########################################
FROM golang:1.17.3-alpine@sha256:55da409cc0fe11df63a7d6962fbefd1321fedc305d9969da636876893e289e2d AS build-env

WORKDIR /app

# entrypoint
RUN apk add --no-cache entr
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

# production image ############################################
FROM alpine:3.15.0@sha256:21a3deaa0d32a8057914f36584b5288d2e5ecc984380bc0118285c70fa8c9300 AS prod
RUN apk add --no-cache ca-certificates

# executable
COPY --from=build-env /go/bin/dind-nurse /bin/dind-nurse
ENTRYPOINT ["dind-nurse"]

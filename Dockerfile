FROM golang:alpine AS builder
COPY * /go/src/github.com/baelish/alive/
RUN apk add git
RUN cd /go/src/github.com/baelish/alive/ && go install .

FROM alpine:latest

ARG BUILD_COMMIT
ARG BUILD_DATE
ARG BUILD_IMAGE
ARG BUILD_VERSION

LABEL org.label-schema.schema-version="1.0"
LABEL org.label-schema.build-date="$BUILD_DATE"
LABEL org.label-schema.name="$BUILD_IMAGE"
LABEL org.label-schema.vcs-url="https://github.com/baelish/alive"
LABEL org.label-schema.vcs-ref="$BUILD_COMMIT"
LABEL org.label-schema.version="$BUILD_VERSION"

RUN mkdir /app
WORKDIR /app
COPY --from=builder /go/bin/alive .
RUN addgroup -g 9001 alive && adduser -h /data -D -u 9001 -G alive alive
USER alive:alive
ENTRYPOINT [ "./alive", "--data-path=/data" , "--static-path=/data/static" ]

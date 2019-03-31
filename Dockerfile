FROM golang:alpine as builder
COPY * /go/src/github.com/baelish/alive/
RUN env
RUN pwd
RUN ls
RUN ls /
RUN find / -name alive.go
RUN cd /go/src/github.com/baelish/alive/ && go install .

FROM alpine:latest
RUN mkdir /app /data
WORKDIR /app
COPY --from=builder /go/bin/alive .
CMD ./alive -b /data/

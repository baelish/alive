FROM golang:alpine as builder
COPY * /go/src/github.com/baelish/alive/
RUN apk add git
RUN go get github.com/gorilla/mux
RUN cd /go/src/github.com/baelish/alive/ && go install .

FROM alpine:latest
RUN mkdir /app /data
WORKDIR /app
COPY --from=builder /go/bin/alive .
CMD ./alive -b /data/

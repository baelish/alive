FROM golang:alpine as builder
RUN env
RUN pwd
RUN ls
RUN ls /
RUN find / -name alive.go
RUN go install .

FROM alpine:latest
RUN mkdir /app /data
WORKDIR /app
COPY --from=builder /go/bin/alive .
CMD ./alive -b /data/

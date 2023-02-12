# Build Geth in a stock Go builder container
FROM golang:1.15-alpine as builder

RUN apk add --no-cache make gcc musl-dev linux-headers git

ADD . /go-virbicoin
RUN cd /go-virbicoin && make gvbc

# Pull Geth into a second stage deploy alpine container
FROM alpine:latest

RUN apk add --no-cache ca-certificates
COPY --from=builder /go-virbicoin/build/bin/gvbc /usr/local/bin/

EXPOSE 8329 8330 28329 28329/udp
ENTRYPOINT ["geth"]

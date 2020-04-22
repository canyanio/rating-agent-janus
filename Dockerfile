FROM golang:1.13.5-alpine3.10 as builder
RUN apk update && apk upgrade && \
    apk add \
    xz-dev \
    musl-dev \
    gcc
RUN mkdir -p /go/src/github.com/canyanio/rating-agent-janus
COPY . /go/src/github.com/canyanio/rating-agent-janus
RUN cd /go/src/github.com/canyanio/rating-agent-janus && env CGO_ENABLED=1 go build

FROM alpine:3.10
RUN apk update && apk upgrade && \
        apk add --no-cache ca-certificates xz
RUN mkdir -p /etc/rating-agent-janus
COPY ./config.yaml /etc/rating-agent-janus
COPY --from=builder /go/src/github.com/canyanio/rating-agent-janus/rating-agent-janus /usr/bin
ENTRYPOINT ["/usr/bin/rating-agent-janus", "--config", "/etc/rating-agent-janus/config.yaml"]

EXPOSE 8000

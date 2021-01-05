FROM golang:1.12-stretch

COPY . /tmp/docker-state-exporter/

WORKDIR /tmp/docker-state-exporter/

RUN env GOOS=linux GOARCH=amd64 go build && ls -l /tmp/docker-state-exporter/

FROM ubuntu:16.04

COPY --from=0 /tmp/docker-state-exporter/docker-state-exporter /usr/local/bin/

CMD ["docker-state-exporter"]
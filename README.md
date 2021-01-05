# docker-state-exporter

Prometheus exporter to exporter container status

## deploy

``` shell
# build image
docker-compose build
# run container
docker-compose up -d
```

## metric

docker_container_state docker container status, ignore some container with prefix-to-skip command arg

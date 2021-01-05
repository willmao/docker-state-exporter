# docker-state-exporter

Prometheus exporter to exporter container status

暴露宿主机上的docker容器状态信息给prometheus。

流程：

- 创建一个docker client，连接到宿主机的docker socket
- 通过client api获取本机所有docker容器
- 过滤掉不关心的容器
- 生成promethues metrics

过滤掉的容器：

- 以k8s_开头的，因为k8s集群中的容器已经被监控了

参数：

- --listen-address 监听端口，默认9901
- --refresh-interval 刷新时间，默认5秒
- --prefix-to-skip 跳过的容器名称，以逗号分隔，默认值k8s_

本项目使用阿里云容器镜像服务自动进行镜像构建，每个release产生一个构建，镜像公网地址: registry.cn-hangzhou.aliyuncs.com/willmao/docker-state-exporter:[版本，如0.9.0]


## 指标

| 名称     | 标签  | 说明  |
|----------|---|---|
| docker_container_state | container_name/exit_code/state  | 容器名称/退出码/当前状态  |
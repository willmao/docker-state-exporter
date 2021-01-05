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
# 使用CentOS 8作为基础镜像
FROM centos:8

# 维护者信息
LABEL maintainer="wangsanchao@Date-02025084"

# 创建应用程序目录
WORKDIR /

# 复制编译好的altervoice应用程序到容器中
COPY altervoice /root

RUN chmod 777 /root/altervoice

# 暴露9000端口（应用程序监听的端口）
EXPOSE 9000

# 设置容器启动时执行的命令
CMD ["/root/altervoice","start"]

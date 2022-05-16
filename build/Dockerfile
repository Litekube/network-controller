FROM centos:7

# docker run -v $ProjectPath:/network-controller
LABEL maintainer="wannazjx@163.com"

# build golang environment
WORKDIR /

RUN yum install wget -y && \
    yum install git -y && \
    yum install tar -y && \
    yum install gcc -y

RUN mkdir -p /usr/local && \
    mkdir -p /home/go && \
    wget --no-check-certificate https://go.dev/dl/go1.17.8.linux-amd64.tar.gz && \
    tar -C /usr/local/ -zxvf go1.17.8.linux-amd64.tar.gz && \
    rm -r go1.17.8.linux-amd64.tar.gz

ENV PATH=/usr/local/go/bin:$PATH
ENV GOPATH /home/go
ENV GOROOT /usr/local/go
ENV GOPROXY goproxy.cn

RUN go version

# build arm-gcc enviroment
RUN yum install glibc.i686 -y && \
    yum install fontconfig.i686 -y && \
    yum install libXrender.i686 -y && \
    yum install libstdc++.i686 -y

COPY ./gcc-linaro-7.5.0-2019.12-i686_arm-linux-gnueabihf.tar.xz /

RUN mkdir -p /usr/local/toolchain/ && \
    xz -d gcc-linaro-7.5.0-2019.12-i686_arm-linux-gnueabihf.tar.xz && \
    tar xvf gcc-linaro-7.5.0-2019.12-i686_arm-linux-gnueabihf.tar -C /usr/local/toolchain/ && \
    rm -r gcc-linaro-7.5.0-2019.12-i686_arm-linux-gnueabihf.tar*

ENV PATH=$PATH:/usr/local/toolchain/gcc-linaro-7.5.0-2019.12-i686_arm-linux-gnueabihf/bin

RUN mkdir -p /network-controller

WORKDIR /network-controller

# RUN chmod +x build/build.sh
CMD ["bin/bash"]
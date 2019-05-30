# IMAGE
FROM golang:1.12.1
# MAINTAINER
MAINTAINER robertzml

# SET FILES
WORKDIR /home/zml/glaucus

ADD . /home/zml/glaucus/

# SET ENVIROMENT
ENV GO111MODULE on
ENV GOPROXY https://goproxy.io

# SET TIMEZONE
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN echo 'Asia/Shanghai' >/etc/timezone

# COMPILE
RUN go build .

EXPOSE 8181

ENTRYPOINT ["./Glaucus"]

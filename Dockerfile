# IMAGE
FROM golang:1.15.5

# SET FILES
WORKDIR /home/zml/glaucus

ADD . /home/zml/glaucus/

# SET ENVIROMENT
ENV GOPROXY https://goproxy.cn

# SET TIMEZONE
# RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
# RUN echo 'Asia/Shanghai' >/etc/timezone

# COMPILE
RUN go build .

ENTRYPOINT ["./Glaucus"]

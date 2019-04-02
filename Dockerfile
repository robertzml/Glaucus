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

# COMPILE
RUN go build

EXPOSE 6540

ENTRYPOINT ["./gorest"]

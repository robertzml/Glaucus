# IMAGE
FROM golang:1.12.1
# MAINTAINER
MAINTAINER robertzml
# SET ENVIROMENT
ENV GO111MODULE on
ENV GOPROXY https://goproxy.io

# SET FILES
RUN mkdir -p /home/zml/glaucus
WORKDIR /home/zml/glaucus
ADD . /home/zml/glaucus

# COMPILE
RUN go build main.go -o glaucus

EXPOSE 8081
ENTRYPOINT  ["./glaucus"]
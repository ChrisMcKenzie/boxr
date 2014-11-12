FROM golang:1.3.3

MAINTAINER Secret Ironman Team

VOLUME /opt/go
ADD . /opt/go/src/github.com/Secret-Ironman/boxr
WORKDIR /opt/go/src/github.com/Secret-Ironman/boxr
ENV GOPATH /opt/go
ENV PATH $PATH:$GOROOT/bin:$GOPATH/bin

RUN make deps && make build-boxr
EXPOSE 3000 3000
CMD bin/boxr s
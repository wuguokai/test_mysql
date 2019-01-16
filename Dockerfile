FROM golang:latest

WORKDIR $GOPATH/src/github.com/wuguokai/test_mysql

ADD . $GOPATH/src/github.com/wuguokai/test_mysql

RUN go build .

EXPOSE 8989

ENTRYPOINT  ["./test_mysql"]
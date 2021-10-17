# Simple compose
FROM golang:1.15

ADD . /src/fib

WORKDIR /src/fib

RUN go mod download

RUN go install ./app/cmd
EXPOSE 8080
ENTRYPOINT /go/bin/cmd
FROM golang:latest

RUN mkdir -p /go/src/app
WORKDIR /go/src/app

ENV GOPATH=/go/src/app/Libraries
ENV CGO_ENABLED=1
ENV GOOS=linux
ENV GOARCH="amd64"
ENV ARGUMENTS=" -f ./helloWord/capibaribe-helloWord-config.yml "

COPY . /go/src/app
COPY docker-entrypoint.sh /docker-entrypoint.sh
RUN chmod +x /docker-entrypoint.sh
RUN cd /go/src/app

#RUN go get github.com/helmutkemper/test_server
#RUN go get -u google.golang.org/grpc
WORKDIR $GOPATH

RUN go build -i -o /tmp/exe /go/src/app/helloWord/client/main.go

EXPOSE 50051:50051

ENTRYPOINT ["/docker-entrypoint.sh"]


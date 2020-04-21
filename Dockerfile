FROM golang:1.13.7-alpine3.11
RUN apk update
RUN apk add --no-cache make nodejs git curl mercurial gcc protobuf protobuf-dev
RUN go get github.com/yudai/gotty
ENV GO111MODULE=on
RUN go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
RUN go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
RUN go get -u github.com/golang/protobuf/protoc-gen-go

WORKDIR /go/src/github.com/logiqai/logiqctl
ADD . .
RUN ./generate_grpc.sh
RUN go build

#
FROM alpine:3.11
EXPOSE 8080
RUN apk update
RUN apk add bash
COPY --from=0 /go/src/github.com/logiqai/logiqctl/logiqctl /bin/logiqctl
COPY --from=0 /go/bin/gotty /bin/gotty
COPY demo.config /root/.logiqctl/config.toml
CMD ["/bin/gotty","-w","/bin/bash"]

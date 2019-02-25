# p2p


Run this in the main dir to compile protobuf
```
protoc -I pkg/service/ --go_out=plugins=grpc:pkg/service/  pkg/service/service.proto
```


Generate reverse-proxy ( reference )

```
protoc -I/usr/local/include -I. \
  -I$GOPATH/src \
  -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
  --grpc-gateway_out=logtostderr=true:. \
  proto/service.proto
  ```
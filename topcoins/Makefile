gentopcoin:
	protoc -I/usr/local/include -I. \
          --grpc-gateway_out=logtostderr=true:./proto \
          --swagger_out=logtostderr=true:./doc \
          --plugin=protoc-gen-grpc-gateway=$(HOME)/go/bin/protoc-gen-grpc-gateway \
          --plugin=protoc-gen-swagger=$(HOME)/go/bin/protoc-gen-swagger\
          --plugin=protoc-gen-go=$(HOME)/go/bin/protoc-gen-go\
           -I$(HOME)/go/src \
           -I$(HOME)/go/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
           --go_out=Mgoogle.golang.org/genproto/googleapis/api/annotations,plugins=grpc:./proto \
          ./api.proto
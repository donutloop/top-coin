genranks:
	protoc -I/usr/local/include -I. \
           -I$(HOME)/go/src \
           --go_out=Mgoogle.golang.org/genproto/googleapis/api/annotations,plugins=grpc:./proto \
          ./api.proto
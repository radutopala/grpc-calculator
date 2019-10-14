all:

generate:
	@echo "Generating code from protos"
	protoc -I . -I$(GOPATH)/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.11.3/third_party/googleapis/ --go_out=plugins=grpc:. --grpc-gateway_out=logtostderr=true:. --swagger_out=logtostderr=true:. ./api/*.proto

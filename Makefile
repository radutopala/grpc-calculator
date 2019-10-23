all:

generate:
	@echo "Generating code from protos"
	protoc -I . -I$(GOPATH)/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.11.3/third_party/googleapis/ --go_out=plugins=grpc:. --grpc-gateway_out=logtostderr=true:. --swagger_out=logtostderr=true:. ./api/*.proto

tests:
	@echo "Running tests"
	go test -bench=. -v ./...

docker-build:
	@echo "Build docker server image"
	docker build -f ./infra/docker/Dockerfile -t radutopala/grpc-calculator:v0.0.1 .
	docker push radutopala/grpc-calculator:v0.0.1

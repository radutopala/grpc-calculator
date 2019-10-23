all:

generate:
	@echo "Generating code from protos"
	protoc -I . -I$(GOPATH)/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.11.3/third_party/googleapis/ --go_out=plugins=grpc:. --grpc-gateway_out=logtostderr=true:. --swagger_out=logtostderr=true:. ./api/*.proto

lint-local:
	@echo "Running linters"
	golangci-lint run -v ./...

lint:
	@echo "Running linters"
	docker run --rm -v $(PWD):/app -w /app golangci/golangci-lint:v1.21.0 golangci-lint run -v ./...

tests:
	@echo "Running tests"
	@mkdir -p artifacts
	go test -race -cover -coverprofile=artifacts/coverage.out -v ./...

coverage: tests
	@echo "Running tests & coverage"
	go tool cover -html=artifacts/coverage.out -o artifacts/coverage.html

benchmark:
	@echo "Running benchmark"
	go test -v -bench=. ./...

docker-build:
	@echo "Build docker server image"
	docker build -f ./infra/docker/Dockerfile -t radutopala/grpc-calculator:v0.0.1 .
	docker push radutopala/grpc-calculator:v0.0.1

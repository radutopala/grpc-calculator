generate:
	protoc -I . --go_out=plugins=grpc:. ./api/*.proto
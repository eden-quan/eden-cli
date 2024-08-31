.PHONY: init
# init and install necessary software
init:
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/golang/mock/mockgen@latest
	go install github.com/go-kratos/kratos/cmd/kratos/v2@latest
	go install github.com/eden/protoc-gen-validate-fx@latest
	go install github.com/eden/protobuf/cmd/protoc-gen-go-fx@latest
	go install github.com/eden/protoc-gen-go-errors-fx@latest
	go install github.com/eden/protoc-gen-openapi-fx@latest
	go install github.com/eden/grpc-gateway/protoc-gen-openapiv2-fx@latest
	go install github.com/eden/protoc-gen-go-http-fx@latest
	go install github.com/eden/protoc-gen-go-grpc-fx@latest
	go install github.com/eden/protoc-gen-go-sql-fx@latest


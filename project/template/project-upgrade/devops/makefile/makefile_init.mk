.PHONY: init
# init and install necessary software
init:
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/golang/mock/mockgen@latest
	go install github.com/go-kratos/kratos/cmd/kratos/v2@latest
	go install gitlab.lainuoniao.cn/rhinobird/backend/protoc-gen-validate-fx.git@latest
	go install gitlab.lainuoniao.cn/rhinobird/backend/protobuf.git/cmd/protoc-gen-go-fx@latest
	go install gitlab.lainuoniao.cn/rhinobird/backend/protoc-gen-go-errors-fx.git@latest
	go install gitlab.lainuoniao.cn/rhinobird/backend/protoc-gen-openapi-fx.git@latest
	go install gitlab.lainuoniao.cn/rhinobird/backend/grpc-gateway.git/protoc-gen-openapiv2-fx@latest
	go install gitlab.lainuoniao.cn/rhinobird/backend/protoc-gen-go-http-fx.git@latest
	go install gitlab.lainuoniao.cn/rhinobird/backend/protoc-gen-go-grpc-fx.git@latest
	go install gitlab.lainuoniao.cn/rhinobird/backend/protoc-gen-go-sql-fx.git@latest

	mv $(shell go env GOPATH)/bin/protoc-gen-validate-fx.git $(shell go env GOPATH)/bin/protoc-gen-validate-fx
	mv $(shell go env GOPATH)/bin/protoc-gen-go-errors-fx.git $(shell go env GOPATH)/bin/protoc-gen-go-errors-fx
	mv $(shell go env GOPATH)/bin/protoc-gen-openapi-fx.git $(shell go env GOPATH)/bin/protoc-gen-openapi-fx
	mv $(shell go env GOPATH)/bin/protoc-gen-go-http-fx.git $(shell go env GOPATH)/bin/protoc-gen-go-http-fx
	mv $(shell go env GOPATH)/bin/protoc-gen-go-grpc-fx.git $(shell go env GOPATH)/bin/protoc-gen-go-grpc-fx
	mv $(shell go env GOPATH)/bin/protoc-gen-go-sql-fx.git $(shell go env GOPATH)/bin/protoc-gen-go-sql-fx

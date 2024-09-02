# generate proto file
COMMON_PROTO=$(shell cd $(PROJECT_PATH) && find api/common -name "*.proto")
API_PROTO=$(shell cd $(PROJECT_PATH) && find api/*$(service)* -name "*.proto")
APP_PROTO=$(shell cd $(PROJECT_PATH) && find app/*$(service)* -name "*.proto")
ALL_PROTO_FILES=$(API_PROTO) $(APP_PROTO) $(COMMON_PROTO)

.PHONY: proto-gen
# generate protobuf with service name : make proto-gen service=service_name
proto-gen:
	@if [ "$(service)" = "" ]; then \
		echo "service is empty, usage: \n\t make proto-gen service=service-name"; \
	else \
		if [ "$(ALL_PROTO_FILES)" != "" ]; then \
			cd $(PROJECT_PATH); \
			protoc \
				--proto_path=. \
				--proto_path=$(GOPATH)/src \
				--proto_path=./third_party \
				--validate-fx_out=paths=source_relative,lang=go:. \
				--go-fx_out=paths=source_relative:. \
				--go-grpc-fx_out=paths=source_relative:. \
				--go-http-fx_out=paths=source_relative:. \
				--go-errors-fx_out=paths=source_relative:. \
				--go-sql-fx_out=paths=source_relative:. \
				--openapiv2-fx_out . \
				--openapiv2-fx_opt logtostderr=true \
				--openapiv2-fx_opt allow_delete_body=true \
				--openapiv2-fx_opt json_names_for_fields=false \
				--openapiv2-fx_opt enums_as_ints=true \
				--openapi-fx_out=fq_schema_naming=true,enum_type=integer,default_response=true:. \
				$(ALL_PROTO_FILES) ; \
		else \
			echo "can't fine proto with service named $(service)"; \
		fi; \
	fi

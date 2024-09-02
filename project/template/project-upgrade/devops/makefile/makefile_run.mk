
SERVICE_PATH := $(shell cd $(PROJECT_PATH) && find ./app -type d -maxdepth 1 -name "*$(service)*" -exec basename {} \;)

.PHONY: run
# run service : example :-->: make run service=grant-service
run:
	@if [ "$(SERVICE_PATH)" = "" ]; then \
		echo "can't find service named $(service)"; \
	else \
		go run ./app/${SERVICE_PATH}/cmd/${SERVICE_PATH} -conf=./app/${SERVICE_PATH}/configs; \
	fi

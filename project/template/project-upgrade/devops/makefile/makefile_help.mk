.DEFAULT_GOAL := help

# show help
help:
	@echo ''
	@echo 'eden project tools'
	@echo 'eden tools provide some utility function for building project'
	@echo '  1.(TODO) make service new={service_name} : initialize an service names with service_name , and the service name with use as the name of go mod'
	@echo '  2.(TODO) make module service={service_name} module={module_name} : create an new module on service'
	@echo '  3.(TODO) make module desc : describe the structure of project'
	@echo '  4. make proto-gen service={service_name} : generate proto files on service specified'
	@echo '  5. make init : initialize the project, setup go env/dependencies, run on first time you got this project'
	@echo ''
	@echo 'Usage: make [Target] [option]'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-_0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[36m%-22s\033[0m %s\n", helpCommand,helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)
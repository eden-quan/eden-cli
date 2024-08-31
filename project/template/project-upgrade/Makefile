GOHOSTOS:=$(shell go env GOHOSTOS)
GOPATH:=$(shell go env GOPATH)
VERSION=$(shell git describe --tags --always)

MAKE_FILE_PATH := $(abspath $(lastword $(MAKEFILE_LIST)))
CURRENT_ABS_PATH=$(shell dirname $(MAKE_FILE_PATH))
CURRENT_PATH=$(shell basename $(CURRENT_ABS_PATH))

PROJECT_PATH=$(shell echo "../../")

jj:
	$(shell $echo("$PROJECT_PATH"))

APP_RELATIVE_PATH=$(shell a=`basename $$PWD` && echo $${a})
ifeq ($(APP_RELATIVE_PATH), $(CURRENT_PATH))
	PROJECT_PATH=./
	APP_RELATIVE_PATH=
endif

ifeq ($(GOHOSTOS), windows)
	#the `find.exe` is different from `find` in bash/shell.
	#to see https://docs.microsoft.com/en-us/windows-server/administration/windows-commands/find.
	#changed to use git-bash.exe to run find cli or other cli friendly, caused of every developer has a Git.
	GIT_BASH= $(subst cmd\,bin\bash.exe,$(dir $(shell where git)))
	COMMON_PROTO_FILES=$(shell $(GIT_BASH) -c "find $(PROJECT_PATH)api/common -name *.proto")
else
endif


include devops/makefile/makefile_init.mk
include devops/makefile/makefile_protoc.mk
include devops/makefile/makefile_run.mk
include devops/makefile/makefile_help.mk

.PHONY: test
# echo test content
test:
	@echo $(CURRENT_PATH)

.PHONY: fmt
## fmt              		: format go files
fmt:
	@echo "--> Formatting go files"
	@go fmt $$(go list ./...)

.PHONY: lint
## lint              		: run lint tools
lint:
	golangci-lint run ./... -v

.PHONY: test
## test             		: run all tests
test:
	go test --race ./... -v

.PHONY: todo
## todo             		: print all todo comments
todo:
	grep --color=always --exclude=GNUmakefile --exclude-dir=.git --exclude-dir=vendor --recursive TODO "$(CURDIR)"

.PHONY: help
help: GNUmakefile
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | sort
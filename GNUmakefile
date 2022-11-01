## fmt      : format go files
## lint     : run lint tools
## test     : run all tests
## tools    : install tools
## todo     : print all todo comments

fmt:
	@echo "--> Formatting go files..."
	@go fmt $$(go list ./...)

lint:
	@echo "==> Checking source code against linters..."
	golangci-lint run ./... -v

test:
	go test --race ./... -v

tools:
	@echo "==> installing required tooling..."
	go install mvdan.cc/gofumpt@latest
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b "$$(go env GOPATH || $$GOPATH)"/bin v1.50.0

todo:
	grep --color=always --exclude=GNUmakefile --exclude-dir=.git --exclude-dir=vendor --recursive TODO "$(CURDIR)"

help: GNUmakefile
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | sort

.PHONY: fmt lint test tools todo help

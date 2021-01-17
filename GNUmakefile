lint:
	golangci-lint run ./... -v

test:
	go test ./... -v

todo:
	grep --color=always --exclude=GNUmakefile --exclude-dir=.git --exclude-dir=vendor --recursive TODO "$(CURDIR)"

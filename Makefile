.PHONY: build
build/cli:
	@echo "+ $@"
	CGO_ENABLED=0 go build -o bin/cli \
        -ldflags "-w -s" \
        github.com/Code-Hex/fast-service/cmd/fast



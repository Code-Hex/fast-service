.PHONY: build
build: build/server build/cli
build/server:
	@echo "+ $@"
	CGO_ENABLED=0 go build -o bin/speedtest-server \
	-ldflags "-w -s" \
        github.com/Code-Hex/fast-service
build/cli:
	@echo "+ $@"
	CGO_ENABLED=0 go build -o bin/speedtest-client \
        -ldflags "-w -s" \
        github.com/Code-Hex/fast-service/cmd/fast



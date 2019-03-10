# fast-service is opensource speedtest service written in Go

<p align="center">
<img src="https://user-images.githubusercontent.com/6500104/54086456-fc5dd100-438c-11e9-9858-2c2e1dacd702.gif"></img>
</p>

## Environment variables

Environment variables is [here](https://github.com/Code-Hex/fast-service/blob/9e3a385f34985237c655efd9aedddbf05ef3ae45/internal/config/config.go#L12-L24)

## How to try this contents

We are necessary [dep](https://github.com/golang/dep#installation) to install dependency packages.

    dep ensure -vendor-only
    make build

### How to run server

    ENV=development ./bin/speedtest-server

### How to run client

If you want to run client on another port (default is 8000), you should rewrite [`cmd/fast/main.go:10`](https://github.com/Code-Hex/fast-service/blob/8c70fbfef8c6efcbd7e6a75e459ec8cf83dde6b5/cmd/fast/main.go#L10) before build client code.

And build it like this.

    ./bin/speedtest-client

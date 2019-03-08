# fast-service is opensource speedtest service written in Go

Deployed onto https://fast.codehex.dev/

## Environment variables

Environment variables is [here](https://github.com/Code-Hex/fast-service/blob/9e3a385f34985237c655efd9aedddbf05ef3ae45/internal/config/config.go#L12-L24)

```sh
export LOG_LEVEL=INFO  # default
export ENV=development # default
export PORT=8000       # default
```

## How to run server

    make build
    ENV=development ./bin/speedtest-server

## How to run client

If you want to run client on development(local), you should rewrite [`cmd/fast/main.go:10`](https://github.com/Code-Hex/fast-service/blob/9e3a385f34985237c655efd9aedddbf05ef3ae45/cmd/fast/main.go#L10) before build client code.

```diff
-const api = "https://fast.codehex.dev"
+const api = "http://localhost:8000" // 8000 represents Port number
```

And build it like this.

    make build/cli
    ./bin/speedtest-client
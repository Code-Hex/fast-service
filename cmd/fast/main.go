package main

import (
	"context"
	"fmt"
	"log"
	"runtime"
)

const api = "https://kouzoh-p-codehex.appspot.com"
const tryCount = 3

var (
	maxConnections = runtime.NumCPU()
	payloadSizes   = []int{
		1562500,  // 1.5625MB
		6250000,  // 6.25MB
		12500000, // 12.5MB
		26214400, // 25MB
	}
)

func main() {
	ctx := context.Background()

	var (
		lastDown  string
		downBytes int64
		lastUp    string
		upBytes   int64
	)
	fmt.Println()
	err := DownloadTest(ctx, func(result *Lap) error {
		lastDown = result.String()
		downBytes = result.Bytes
		fmt.Printf("    %s, size: %d ↓ - %s bps, size: %d ↑\r", lastDown, downBytes, "", 0)
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	err = UploadTest(ctx, func(result *Lap) error {
		lastUp = result.String()
		upBytes = result.Bytes
		fmt.Printf("    %s, size: %d ↓ - %s, size: %d ↑\r", lastDown, downBytes, lastUp, result.Bytes)
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("    %s, size: %d ↓ - %s, size: %d ↑\n", lastDown, downBytes, lastUp, upBytes)
}

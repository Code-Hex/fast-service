package main

import (
	"context"
	"time"

	"golang.org/x/sync/errgroup"
)

var DownloadTimeout = 15 * time.Second

const downloadURL = api + "/download"

func DownloadTest(ctx context.Context, cb IntervalCallback) error {
	eg, ctx := errgroup.WithContext(ctx)
	ctx, cancel := context.WithTimeout(ctx, DownloadTimeout)
	defer cancel()

	r := newRecord(time.Now(), maxConnections)

	go func() {
		for {
			select {
			case lap := <-r.Lap():
				cb(&lap)
			case <-ctx.Done():
				return
			}
		}
	}()

	semaphore := make(chan struct{}, maxConnections)
loop:
	for i := 0; i < tryCount; i++ {
		for _, size := range payloadSizes {
			select {
			case <-ctx.Done():
				break loop
			case semaphore <- struct{}{}:
				time.Sleep(250 * time.Millisecond)
			}
			eg.Go(func() error {
				defer func() { <-semaphore }()
				if err := r.download(ctx, downloadURL, size); err != nil {
					return err
				}
				return nil
			})
		}
	}
	// waiting
	select {
	case <-ctx.Done():
	case semaphore <- struct{}{}:
		cancel()
	}
	return errorCheck(eg.Wait())
}

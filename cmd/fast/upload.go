package main

import (
	"context"
	"net/url"
	"time"

	"golang.org/x/sync/errgroup"
)

var UploadTimeout = 15 * time.Second

const uploadURL = api + "/upload"

type IntervalCallback func(current *Lap) error

func UploadTest(ctx context.Context, cb IntervalCallback) error {
	eg, ctx := errgroup.WithContext(ctx)
	ctx, cancel := context.WithTimeout(ctx, UploadTimeout)
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
				if err := r.upload(ctx, uploadURL, size); err != nil {
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

func errorCheck(err error) error {
	if err == context.Canceled || err == context.DeadlineExceeded {
		return nil
	}
	if v, ok := err.(*url.Error); ok {
		err = v.Err
		return errorCheck(err)
	}
	return err
}

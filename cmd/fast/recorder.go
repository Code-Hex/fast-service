package main

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sync/atomic"
	"time"
)

type recorder struct {
	byteLen int64
	start   time.Time
	lapch   chan Lap
}

func newRecorder(start time.Time, cpun int) *recorder {
	return &recorder{
		start: start,
		lapch: make(chan Lap, cpun),
	}
}

func (r *recorder) Lap() <-chan Lap {
	return r.lapch
}

func (r *recorder) download(ctx context.Context, url string, size int) error {
	url = fmt.Sprintf("%s?size=%d", url, size)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		return err
	}

	// status check
	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
	}

	// start measure
	proxy := r.newRecordProxy(ctx, resp.Body)
	defer proxy.Close()

	if _, err := io.Copy(ioutil.Discard, proxy); err != nil {
		return err
	}
	return nil
}

func (r *recorder) upload(ctx context.Context, url string, size int) error {
	// start measure
	proxy := r.newRecordProxy(ctx, rand.Reader)
	defer proxy.Done()
	req, err := http.NewRequest("POST", url, proxy)
	if err != nil {
		return err
	}
	req.ContentLength = int64(size)
	req.Header.Set("Content-Type", "application/octet-stream")

	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		return err
	}

	// status check
	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
	}

	return nil
}

type recordProxy struct {
	context.Context
	io.Reader
	*recorder
	done chan struct{}
}

func (r *recorder) newRecordProxy(ctx context.Context, reader io.Reader) *recordProxy {
	rp := &recordProxy{
		Context:  ctx,
		Reader:   reader,
		recorder: r,
		done:     make(chan struct{}),
	}
	go rp.Watch(r.lapch)
	return rp
}

func (r *recordProxy) Done() { close(r.done) }

type newrecord struct {
	Bytes     int64
	Bps       float64
	PrettyBps float64
	Units     string
}

func (r *recordProxy) Watch(send chan<- Lap) {
	t := time.NewTicker(150 * time.Millisecond)
	for {
		select {
		case <-t.C:
			byteLen := atomic.LoadInt64(&r.byteLen)
			delta := time.Now().Sub(r.start).Seconds()
			send <- newLap(byteLen, delta)
		case <-r.done:
			return
		}
	}
}

func (r *recordProxy) Read(p []byte) (n int, err error) {
	select {
	case <-r.Context.Done():
		return 0, nil
	default:
	}
	n, err = r.Reader.Read(p)
	if err != nil {
		return 0, err
	}
	atomic.AddInt64(&r.byteLen, int64(n))
	return
}

// Close the reader when it implements io.Closer
func (r *recordProxy) Close() error {
	if closer, ok := r.Reader.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}

package main

import (
	"fmt"

	humanize "github.com/dustin/go-humanize"
)

// Lap represents lap used in stopwatch.
type Lap struct {
	Bytes     int64
	Bps       float64
	PrettyBps float64
	Units     string

	delta float64
}

func newLap(byteLen int64, delta float64) Lap {
	var bytes float64
	if delta > 0 {
		bytes = float64(byteLen) / delta
	}
	bps := bytes * 8
	prettyBps, unit := humanize.ComputeSI(bps)
	return Lap{
		Bytes:     byteLen,
		Bps:       bps,
		PrettyBps: prettyBps,
		Units:     unit + "bps",

		delta: delta,
	}
}

func (l *Lap) String() string {
	return fmt.Sprintf("%7.2f %s", l.PrettyBps, l.Units)
}

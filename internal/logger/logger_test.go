package logger

import (
	"testing"

	"go.uber.org/zap/zapcore"
)

func TestLogLevel(t *testing.T) {
	cases := []struct {
		level   string
		success bool
		want    zapcore.Level
	}{
		{
			"info",
			true,
			zapcore.InfoLevel,
		},
		{
			"DEBUG",
			true,
			zapcore.DebugLevel,
		},
		{
			"Error",
			true,
			zapcore.ErrorLevel,
		},
		{
			"FATAL", // not supported (debug or info is enough)
			false,
			zapcore.Level(0),
		},
	}

	for _, tc := range cases {
		got, err := logLevel(tc.level)
		if err != nil {
			if tc.success {
				t.Fatalf("expect to success: %s", err)
			}
			continue
		}

		if !tc.success {
			t.Fatal("expect to be failed")
		}

		if got != tc.want {
			t.Fatalf("got %v, want %v", got, tc.want)
		}
	}
}

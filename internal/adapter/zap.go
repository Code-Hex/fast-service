package adapter

import (
	"net/http"
	"time"

	"github.com/Code-Hex/fast-service/internal/logger"
	"github.com/rs/xid"
	"go.uber.org/zap"
)

type delegator struct {
	http.ResponseWriter
	Status int
}

// ZapAdapter returns Adapter which adapt logging middleware.
func ZapAdapter(l *zap.Logger) Adapter {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			ctx := logger.ToContext(
				r.Context(),
				l.With(zap.Stringer("request_id", xid.New())),
			)

			d := &delegator{ResponseWriter: w}
			next.ServeHTTP(d, r.WithContext(ctx))

			logger.Info(ctx, "request",
				zap.String("host", r.Host),
				zap.String("path", r.URL.Path),
				zap.Int("status", d.Status),
				zap.Duration("duration", time.Now().Sub(start)),
				zap.String("method", r.Method),
				zap.String("user_agent", r.UserAgent()),
			)
		})
	}
}

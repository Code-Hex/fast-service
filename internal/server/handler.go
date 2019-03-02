package server

import (
	"net/http"

	"github.com/Code-Hex/fast-service/internal/adapter"
	"go.uber.org/zap"
)

// Mux manages handlers.
type Mux struct {
	mux      *http.ServeMux
	adapters []adapter.Adapter
}

// NewMux creates new Mux http.Handler.
func NewMux(logger *zap.Logger) *Mux {
	return &Mux{
		mux: http.NewServeMux(),
		// if set adapter like []Adapter{A(), B()}
		// we will access A() -> B() -> main -> B() -> A()
		adapters: []adapter.Adapter{
			adapter.ZapAdapter(logger),
		},
	}
}

// Handle registers the handler for the given pattern.
func (m *Mux) Handle(pattern string, h http.Handler) {
	m.mux.Handle(pattern, adapter.Adapt(h, m.adapters...))
}

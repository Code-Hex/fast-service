package adapter

// https://medium.com/@matryer/writing-middleware-in-golang-and-how-go-makes-it-so-much-fun-4375c1246e81
import "net/http"

// Adapter represents middleware between request and main handler
type Adapter func(http.Handler) http.Handler

// Adapt adapts datapters to giving handler.
func Adapt(h http.Handler, adapters ...Adapter) http.Handler {
	// To process from left to right, iterate from the last one.
	for i := len(adapters) - 1; i >= 0; i-- {
		h = adapters[i](h)
	}
	return h
}

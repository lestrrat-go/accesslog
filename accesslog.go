package accesslog

import (
	"log/slog"
	"net/http"
)

// Middleware is the main object that you interact with. Despite its name
// it's actually a Builder object, which needs to be applied to a http.Handler
// by calling `Wrap()`
//
// For frameworks such as CHI that require a function that returns a http.Handler,
// pass it the reference to the `Wrap()` method.
type Middleware struct {
	clock          Clock
	collector      Collector
	logger         *slog.Logger
	logLevel       slog.Level
	recordResponse bool
	rwbuilder      ResponseWriterBuilder
}

// Collector is an object that collects attributes from the request and response
// and returns them as a slice of slog.Attr objects to be logged.
//
// If you do not like the default collector that is provided by the library,
// you can always implement your own. If you just want to extend the standard
// collector, you can embed the `Standard()` collector to your collector, and
// add your own attributes to the return values
type Collector interface {
	Collect(ResponseWriter, *http.Request) []slog.Attr
}

type handler struct {
	next           http.Handler
	clock          Clock
	collector      Collector
	logger         *slog.Logger
	logLevel       slog.Level
	recordResponse bool
	rwbuilder      ResponseWriterBuilder
}

// New creates a new Middleware object. You can further customize the object
// by calling its methods.
//
// By default the `Standard()` collector is added to the list of collectors,
// and the log level is set to `slog.LevelInfo`
func New() *Middleware {
	return &Middleware{
		clock:     SystemClock{},
		collector: Standard(),
		logLevel:  slog.LevelInfo,
		logger:    slog.Default(),
		rwbuilder: DefaultResponseWriterBuilder(),
	}
}

// Clock sets the Clock object to be used by the Middleware object. By default
// `SystemClock` is used.
func (al *Middleware) Clock(clock Clock) *Middleware {
	al.clock = clock
	return al
}

// Collector sets the Collector object to be used by the Middleware object. By
// default `Standard()` is used.
func (al *Middleware) Collector(collector Collector) *Middleware {
	al.collector = collector
	return al
}

// Logger sets the slog.Logger object to be used by the Middleware object. By
// default `slog.Default()` is used.
func (al *Middleware) Logger(logger *slog.Logger) *Middleware {
	al.logger = logger
	return al
}

// LogLevel sets the log level to be used by the Middleware object. By default
// `slog.LevelInfo` is used.
func (al *Middleware) LogLevel(level slog.Level) *Middleware {
	al.logLevel = level
	return al
}

// RecordResponse sets whether the response body should be recorded. By default
// this is set to false for performance reasons.
func (al *Middleware) RecordResponse(b bool) *Middleware {
	al.recordResponse = b
	return al
}

// Wrap returns a http.Handler that wraps the provided `next“ http.Handler object.
func (al *Middleware) Wrap(next http.Handler) http.Handler {
	return handler{
		next:           next,
		clock:          al.clock,
		collector:      al.collector,
		logger:         al.logger,
		logLevel:       al.logLevel,
		recordResponse: al.recordResponse,
		rwbuilder:      al.rwbuilder,
	}
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rw := h.rwbuilder.Wrap(w, r, h.recordResponse)
	h.next.ServeHTTP(rw, r)
	rw.End()
	h.process(rw, r)
}

func (m handler) process(rw ResponseWriter, r *http.Request) {
	attrs := m.collector.Collect(rw, r)
	m.logger.LogAttrs(r.Context(), m.logLevel, "access", attrs...)
}

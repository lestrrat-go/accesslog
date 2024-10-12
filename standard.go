package accesslog

import (
	"log/slog"
	"net/http"
)

// Standard collects the "standard" set of attributes.
func Standard() Collector {
	return standard{}
}

type standard struct{}

func (standard) Collect(rw ResponseWriter, r *http.Request) []slog.Attr {
	return []slog.Attr{
		slog.String("remote_addr", r.RemoteAddr),
		slog.String("http_method", r.Method),
		slog.String("path", r.URL.Path),
		slog.Int("status", rw.Status()),
		slog.Int64("body_bytes_sent", rw.BytesWritten()),
		slog.String("http_referer", r.Referer()),
		slog.String("http_user_agent", r.UserAgent()),
	}
}

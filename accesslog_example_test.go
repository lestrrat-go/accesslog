package accesslog_test

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"time"

	"github.com/lestrrat-go/accesslog"
)

func Example() {
	al := accesslog.New().
		// Set the clock to a static time to force duration=0 for testing
		Clock(accesslog.StaticClock(time.Time{})).
		Logger(
			slog.New(
				slog.NewJSONHandler(os.Stdout,
					&slog.HandlerOptions{
						ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
							switch a.Key {
							case slog.TimeKey:
								// replace time to get static output for testing
								return slog.Time(slog.TimeKey, time.Time{})
							case "remote_addr":
								// replace value to get static output for testing
								return slog.String("remote_addr", "127.0.0.1:99999")
							}
							return a
						},
					},
				),
			),
		)

	srv := httptest.NewServer(al.Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("Hello, World!"))
	})))
	defer srv.Close()

	_, err := http.Get(srv.URL)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// OUTPUT:
	// {"time":"0001-01-01T00:00:00Z","level":"INFO","msg":"access","remote_addr":"127.0.0.1:99999","http_method":"GET","path":"/","status":200,"body_bytes_sent":13,"http_referer":"","http_user_agent":"Go-http-client/1.1"}
}

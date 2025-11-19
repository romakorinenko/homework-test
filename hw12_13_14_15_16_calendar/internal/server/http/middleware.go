package internalhttp

import (
	"fmt"
	"net/http"
	"time"

	"github.com/romakorinenko/homework-test/hw12_13_14_15_calendar/internal/app"
)

func loggingMiddleware(next http.Handler, log app.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Info(fmt.Sprintf("%s  %s  %s  %s  %s  %v",
			r.RemoteAddr,
			r.Method,
			r.URL.Path,
			r.Proto,
			r.UserAgent(),
			time.Since(start)))
	}
}

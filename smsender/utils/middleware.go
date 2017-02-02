package utils

import (
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
)

func Logger(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	start := time.Now()
	path := r.URL.Path

	next(w, r)

	end := time.Now()

	log.Printf(
		"%s %s %s %13v",
		end.Format("2006/01/02 - 15:04:05"),
		r.Method,
		path,
		end.Sub(start),
	)
}

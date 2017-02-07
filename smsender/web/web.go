package web

import (
	"net/http"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/minchao/smsender/smsender"
	config "github.com/spf13/viper"
	"github.com/urfave/negroni"
)

func InitWeb(sender *smsender.Sender) {
	log.Debug("web.InitWeb")

	if config.GetBool("http.web.enable") {
		router := sender.HTTPRouter

		router.PathPrefix("/static/").
			Handler(staticHandler(http.StripPrefix("/static/", http.FileServer(http.Dir("./webroot/public/")))))

		n := negroni.New(negroni.Wrap(http.HandlerFunc(root)))

		router.Handle("/", n).Methods("GET")
		router.Handle("/{anything:.*}", n).Methods("GET")
	}
}

func root(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache, max-age=31556926, public")
	http.ServeFile(w, r, "./webroot/index.html")
}

func staticHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "max-age=31556926, public")
		if strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(w, r)
			return
		}
		handler.ServeHTTP(w, r)
	})
}

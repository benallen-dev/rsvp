package web

import (
	"net/http"
	"rsvp/store"
	"rsvp/web/handlers"

	"github.com/charmbracelet/log"
)

func NewMux(s *store.Store) *http.ServeMux {
	mux := http.NewServeMux()

	fileHandler := http.FileServer(http.Dir("./public"))

	mux.Handle("GET /static/", http.StripPrefix("/static/", fileHandler))

	mux.HandleFunc("GET /overview", handlers.GetOverview(s)) // For BTS use

	mux.HandleFunc("POST /rsvp", handlers.PostRsvp(s)) // Submits the form

	mux.HandleFunc("GET /day/rsvp", handlers.GetRsvp())
	mux.HandleFunc("GET /evening/rsvp", handlers.GetRsvp())
	mux.HandleFunc("GET /day", handlers.GetHome())
	mux.HandleFunc("GET /evening", handlers.GetHome())

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		p := r.URL.Path

		if p == "/" {
			log.Debug("Got root request, serving home", "path", p)
			handlers.GetHome()(w,r)
		} else {
			log.Debug("Got root request, serving file", "path", p)
			fileHandler.ServeHTTP(w,r)
		}
	})

	return mux
}

func noCacheMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
		next.ServeHTTP(w, r)
	})
}

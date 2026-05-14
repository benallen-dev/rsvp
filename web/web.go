package web

import (
	"net/http"
	"rsvp/store"
	"rsvp/web/handlers"
)

func NewMux(s *store.Store) *http.ServeMux {
	mux := http.NewServeMux()
	
	mux.HandleFunc("GET /", handlers.GetHome)

	mux.HandleFunc("GET /overview", handlers.GetOverview(s)) // For BTS use

	mux.HandleFunc("GET /rsvp", handlers.GetRsvp(s)) // Gets the form
	mux.HandleFunc("POST /rsvp", handlers.PostRsvp(s)) // Submits the form

	// Catch-all for things like favicon
	// fileHandler := http.FileServer(http.Dir("./public"))
	// mux.Handle("GET /", noCacheMiddleware(fileHandler))
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./public"))))
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

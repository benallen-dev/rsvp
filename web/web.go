package web

import (
	"net/http"
	"rsvp/store"
)

func NewMux(store *store.Store) *http.ServeMux {
	mux := http.NewServeMux()

	fileHandler := http.FileServer(http.Dir("./public"))
	mux.Handle("GET /", noCacheMiddleware(fileHandler))

	mux.HandleFunc("GET /rsvp", getRsvpIndex(store))
	mux.HandleFunc("GET /rsvp/{inviteId}", getRsvp(store))
	mux.HandleFunc("POST /rsvp/{inviteId}", postRsvp(store))

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

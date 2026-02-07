package web

import (
	"net/http"
)

func NewMux() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /rsvp/{inviteId}", getRsvp)
	mux.HandleFunc("POST /rsvp/{inviteId}", postRsvp)

	return mux
}

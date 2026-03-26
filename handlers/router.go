package handlers

import (
	"net/http"

	"rsvp/storage"
)

func NewMux(s *storage.Store) *http.ServeMux {
	mux := http.NewServeMux()

	// Create middleware chain for store injection
	storeInject := storeMiddleware(s)

	mux.HandleFunc("GET /", noCacheHandler(getHome))

	mux.HandleFunc("GET /invite", wrapHandler(storeInject, getInvites))
	mux.HandleFunc("GET /invite/{inviteId}", wrapHandler(storeInject, getInvite))

	mux.HandleFunc("GET /rsvp", wrapHandler(storeInject, getRsvpIndex))
	mux.HandleFunc("GET /rsvp/{inviteId}", wrapHandler(storeInject, getRsvp))
	mux.HandleFunc("POST /rsvp/{inviteId}", wrapHandler(storeInject, postRsvp))

	mux.HandleFunc("GET /respond", wrapHandler(storeInject, getRespond))
	mux.HandleFunc("GET /respond/form", wrapHandler(storeInject, getRespondFormFields))

	return mux
}

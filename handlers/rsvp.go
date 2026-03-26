package handlers

import (
	"net/http"

	"github.com/charmbracelet/log"
)

func getRsvpIndex(w http.ResponseWriter, r *http.Request) {
	log.Info("GET /rsvp")

	s, err := StoreFromContext(r.Context())
	if err != nil {
		log.Error("store not available", "err", err)
		http.Error(w, "Store unavailable", http.StatusInternalServerError)
		return
	}

	_ = s // TODO: Use store to fetch RSVPs
	w.Write([]byte("TODO: Render RSVP template"))
}

func getRsvp(w http.ResponseWriter, r *http.Request) {
	inviteId := r.PathValue("inviteId")
	log.Info("GET /rsvp", "inviteId", inviteId)

	s, err := StoreFromContext(r.Context())
	if err != nil {
		log.Error("store not available", "err", err)
		http.Error(w, "Store unavailable", http.StatusInternalServerError)
		return
	}

	_ = s // TODO: Use store to fetch RSVP data for invite
	w.Write([]byte("TODO: Render RSVP template"))
}

package handlers

import (
	"encoding/json"
	"net/http"
	"rsvp/invite"
	"rsvp/store"

	"github.com/charmbracelet/log"
)

// Returns the form for a given rsvp
func GetRsvp(s *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// inviteId := r.FormValue("invite-id")
		inviteId := r.PathValue("inviteId")

		log.Info("GET /rsvp", "inviteId", inviteId)

		inv, err := s.ReadInvite(inviteId)
		if err != nil {
			http.Error(w, "Could not fetch invite for id "+inviteId, 400)
		}

		log.Infof("%v", inv)

		w.Write([]byte("TODO: Render RSVP template for " + inv.Name + " <small>(" + inviteId + ")</small>"))
	}
}

func PostRsvp(s *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		inviteId := r.PathValue("inviteId")
		log.Info("POST /rsvp", "id", inviteId)

		decoder := json.NewDecoder(r.Body)

		b, err := invite.NewRsvp(inviteId)
		if err != nil {
			log.Error("could not create rsvp for inviteId" + inviteId)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("could not create rsvp for inviteId " + inviteId))
			return
		}

		err = decoder.Decode(&b)
		if err != nil {
			log.Error("could not decode RSVP", "id", inviteId)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		log.Warn("TODO: Actually do something with the RSVP")

		log.Info(b)

		w.Write([]byte("TODO: Render received template\n"))
	}
}

package web

import (
	"encoding/json"
	"net/http"
	"rsvp/invite"

	"github.com/charmbracelet/log"
)

func getRsvpIndex(w http.ResponseWriter, r *http.Request) {
	log.Info("GET /rsvp",)

	w.Write([]byte("TODO: Render RSVP template"))
}

func getRsvp(w http.ResponseWriter, r *http.Request) {
	inviteId := r.PathValue("inviteId")
	log.Info("GET /rsvp", "inviteId", inviteId)


	w.Write([]byte("TODO: Render RSVP template"))
}

func postRsvp(w http.ResponseWriter, r *http.Request) {
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

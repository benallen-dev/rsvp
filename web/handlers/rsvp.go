package handlers

import (
//	"encoding/json"
	"html/template"
	"net/http"
	
//	"rsvp/invite"
	"rsvp/store"

	"github.com/charmbracelet/log"
)

var rsvpTemplate *template.Template

type rsvpTemplateData struct {
	Mode string
}

func init() {
	// Silent fail on init - templates might not exist yet during startup
	var err error
	rsvpTemplate, err = template.New("rsvp-tpl").Funcs(funcMap).ParseFiles(
		"web/templates/base.html",
		"web/templates/rsvp/rsvp.html",
		"web/templates/rsvp/_form-day.html",
		"web/templates/rsvp/_form-evening.html",
	)
	if err != nil {
		log.Warn("RSVP templates not yet available", "err", err)
	}
}
// Returns the form for a given rsvp
func GetRsvp() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("rsvp")

		routeType := getRouteType(r.URL)
		log.Infof("Route type: %s",routeType)

		if routeType != "day" && routeType != "evening" {
			log.Info("redirecting")
			http.Redirect(w, r, "/", http.StatusFound)
			return;
		}

		log.Info("Rendering template")

	
		// Get invites and pass them to template
		var tplData rsvpTemplateData
		log.Infof("Route type: %s",routeType)
		tplData.Mode = routeType

		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")

		if err := rsvpTemplate.ExecuteTemplate(w, "base", tplData); err != nil {
			log.Error("Template execution failed", "err", err)
		}
		return;
	}
}

func PostRsvp(s *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("POST /rsvp")

		// decoder := json.NewDecoder(r.Body)

		// b, err := invite.NewRsvp(inviteId)
		// if err != nil {
		// 	log.Error("could not create rsvp for inviteId" + inviteId)
		// 	w.WriteHeader(http.StatusBadRequest)
		// 	w.Write([]byte("could not create rsvp for inviteId " + inviteId))
		// 	return
		// }

		// err = decoder.Decode(&b)
		// if err != nil {
		// 	log.Error("could not decode RSVP", "id", inviteId)
		// 	w.WriteHeader(http.StatusBadRequest)
		// 	return
		// }

		log.Warn("TODO: Actually do something with the RSVP")

		// log.Info(b)

		w.Write([]byte("TODO: Render received template\n"))
	}
}

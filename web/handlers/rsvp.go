package handlers

import (
	"html/template"
	"net/http"
	"time"
	
	"rsvp/invite"
	"rsvp/store"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
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
		log.Debugf("Route type: %s",routeType)

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

		err := r.ParseForm()
		if err != nil {
			log.Error("could not parse form", "err", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		name := r.FormValue("name")
		rsvpType := r.FormValue("form-type")
		attendingDay := r.FormValue("attending_day") == "on"
		attendingEvening := r.FormValue("attending_evening") == "on"
		dietNotes := r.FormValue("diet_notes")
		message := r.FormValue("message")

		if name == "" {
			log.Error("name required")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Name is required"))
			return
		}

		rsvp := &invite.RSVP{
			Id:               uuid.New(),
			Timestamp:        time.Now(),
			RsvpType:         rsvpType,
			Name:             name,
			AttendingDay:     attendingDay,
			AttendingEvening: attendingEvening,
			DietNotes:        dietNotes,
			Message:          message,
		}

		log.Infof("Creating RSVP: %v", rsvp)

		_, err = s.CreateRSVP(rsvp)
		if err != nil {
			log.Error("could not create rsvp", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed to save RSVP"))
			return
		}

		log.Info("RSVP created", "name", name, "type", rsvpType)
		w.Write([]byte("Thanks for your RSVP!"))
	}
}

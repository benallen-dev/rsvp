package handlers

import (
	"encoding/base64"
	"encoding/json"
	"errors"
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
	Mode         string
	HasRsvp      bool
	ExistingRsvp *invite.RSVP
}

// returns pointer if ok, nil on ErrNoCookie, error on corruption
func decodeCookie(r *http.Request) (*invite.RSVP, error) {
	cookieData, err := r.Cookie("formdata")
	if err == http.ErrNoCookie {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	cookieValue, err := base64.StdEncoding.DecodeString(cookieData.Value)
	if err != nil {
		return nil, errors.New("Could not decode cookie")
	}

	var cookieRsvp invite.RSVP
	err = json.Unmarshal(cookieValue, &cookieRsvp)
	if err != nil {
		return nil, err
	}

	return &cookieRsvp, nil
}

func init() {
	// Silent fail on init - templates might not exist yet during startup
	var err error
	rsvpTemplate, err = template.New("rsvp-tpl").Funcs(funcMap).ParseFiles(
		"web/templates/base.html",
		"web/templates/rsvp/rsvp.html",
		"web/templates/rsvp/_rsvp-display.html",
		"web/templates/rsvp/_form-day.html",
		"web/templates/rsvp/_form-evening.html",
		"web/templates/rsvp/_form-submit.html",
	)
	if err != nil {
		log.Warn("RSVP templates not yet available", "err", err)
	}
}

// Returns the form for a given rsvp
func GetRsvp() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		routeType := getRouteType(r.URL)
		log.Debugf("[RSVP] Route type: %s", routeType)

		if routeType != "day" && routeType != "evening" {
			log.Info("redirecting")
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		cookieRsvp, err := decodeCookie(r)
		if err != nil {
			log.Warn(err.Error())
			http.SetCookie(w, &http.Cookie{
				Name:   "formdata",
				Value:  "",
				MaxAge: -1,
			})
		}

		// Get invites and pass them to template
		var tplData rsvpTemplateData
		log.Infof("Route type: %s", routeType)
		tplData.Mode = routeType
		tplData.HasRsvp = cookieRsvp != nil // kinda gross doing a null check in go but whatever
		tplData.ExistingRsvp = cookieRsvp

		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")

		if err := rsvpTemplate.ExecuteTemplate(w, "base", tplData); err != nil {
			log.Error("Template execution failed", "err", err)
		}
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
		attendingCeremony := r.FormValue("attending_ceremony") == "on"
		attendingReception := r.FormValue("attending_reception") == "on"
		attendingDinner := r.FormValue("attending_dinner") == "on"
		attendingParty := r.FormValue("attending_party") == "on"
		dietNotes := r.FormValue("diet_notes")
		message := r.FormValue("message")

		if name == "" {
			log.Error("name required")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Name is required"))
			return
		}

		rsvp := &invite.RSVP{
			Id:                 uuid.New(),
			Timestamp:          time.Now(),
			Type:               rsvpType,
			Name:               name,
			DietNotes:          dietNotes,
			Message:            message,
			AttendingCeremony:  attendingCeremony,
			AttendingReception: attendingReception,
			AttendingDinner:    attendingDinner,
			AttendingParty:     attendingParty,
		}

		log.Infof("Creating RSVP: %v", rsvp)

		dbRsvp, err := s.CreateRSVP(rsvp)
		if err != nil {
			log.Error("could not create rsvp", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed to save RSVP"))
			return
		}

		// stringify some json of that sweet rsvp
		rsvpBytes, err := json.Marshal(dbRsvp)
		if err != nil {
			log.Error("could not serialise rsvp", "err", err)
		}

		rsvpBytesPretty, err := json.MarshalIndent(rsvp, "", "  ")
		if err != nil {
			log.Error("could not serialise rsvp", "err", err)
		}

		http.SetCookie(w, &http.Cookie{
			Name:        "formdata",
			Value:       base64.StdEncoding.EncodeToString(rsvpBytes),
			Quoted:      false,
			Path:        "/",
			Domain:      "",
			Expires:     time.Date(2026, time.August, 28, 12, 0, 0, 0, time.UTC),
			Secure:      true,
			HttpOnly:    false,
			SameSite:    http.SameSiteLaxMode,
			Partitioned: false,
		})

		log.Info("RSVP created", "name", name, "type", rsvpType)

		if err := rsvpTemplate.ExecuteTemplate(w, "form-submit", dbRsvp); err != nil {
			log.Error("Template execution failed", "err", err)
			w.Write([]byte("Thanks for your RSVP!<br /><pre>" + string(rsvpBytesPretty) + "</pre>"))
		}
	}
}

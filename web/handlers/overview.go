package handlers

import (
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/charmbracelet/log"

	"rsvp/config"
	"rsvp/invite"
	"rsvp/store"
)

var overviewTemplate *template.Template

type OverviewData struct {
	Invites     []*invite.InviteWithRSVPs
	Stringified string
}

func init() {
	// Silent fail on init - templates might not exist yet during startup
	var err error
	overviewTemplate, err = template.ParseFiles(
		"web/templates/base.html",
		"web/templates/overview/overview.html",
	)
	if err != nil {
		log.Debug("Home templates not yet available", "err", err)
	}
}

func GetOverview(s *store.Store) http.HandlerFunc {
	cfg := config.Current

	return func(w http.ResponseWriter, r *http.Request) {
		if config.Current.AuthEnabled == true {
			user, pass, ok := r.BasicAuth()
			if !ok || user != cfg.Users.Admin.Username || pass != cfg.Users.Admin.Password {
				w.Header().Set("WWW-Authenticate", `Basic realm="restricted"`)
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Unauthorized"))
				return
			}
		}

		if overviewTemplate == nil {
			// I think we can remove this for prod? Idk.
			var err error
			overviewTemplate, err = template.ParseFiles(
				"web/templates/base.html",
				"web/templates/overview/overview.html",
			)
			if err != nil {
				log.Warn("Template reload failed", "err", err)
				http.Error(w, "Template error", http.StatusInternalServerError)
				return
			}
			log.Info("Home templates reloaded (dev mode)")
		}

		data, err := s.ReadAllInvitesWithRSVPs()
		if err != nil {
			log.Error(err.Error())
			http.Error(w, "Data error", http.StatusInternalServerError)
			return
		}

		jsonBytes, err := json.Marshal(data)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, "Data error", http.StatusInternalServerError)
			return
		}

		templateData := OverviewData{
			Invites: data,
			Stringified: string(jsonBytes),
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")

		if err := overviewTemplate.ExecuteTemplate(w, "base", templateData); err != nil {
			log.Error("Template execution failed", "err", err)
		}
	}
}

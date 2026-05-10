package handlers

import (
	"html/template"
	"net/http"

	"github.com/charmbracelet/log"

	"rsvp/config"
	"rsvp/store"
)

var overviewTemplate *template.Template

const USERNAME = "ben"
const PASSWORD = "foo"

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

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")

		if err := overviewTemplate.ExecuteTemplate(w, "base", nil); err != nil {
			log.Error("Template execution failed", "err", err)
		}
	}
}

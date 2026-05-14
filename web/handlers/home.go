package handlers

import (
	"html/template"
	"net/http"

	"github.com/charmbracelet/log"

	"rsvp/invite"
	"rsvp/store"
)

var homeTemplate *template.Template

type homeTemplateData struct {
	Invites []*invite.Invite
}

func init() {
	// Silent fail on init - templates might not exist yet during startup
	var err error
	homeTemplate, err = template.ParseFiles(
		"web/templates/base.html",
		"web/templates/home/home.html",
		"web/templates/home/_hero.html",
		"web/templates/home/_timeline.html",
		"web/templates/home/_directions.html",
		"web/templates/home/_dresscode.html",
		"web/templates/home/_contact.html",
		"web/templates/home/_rsvp-cta.html",
		"web/templates/home/_rsvp-frame.html",
		"web/templates/home/_sticky-footer.html",
	)
	if err != nil {
		log.Debug("Home templates not yet available", "err", err)
	}
}



func GetHome(s *store.Store) http.HandlerFunc {
	//cfg := config.Current

	return func(w http.ResponseWriter, r *http.Request) {
		// Reload on error to allow dev editing without restart
		if homeTemplate == nil {
			var err error
			homeTemplate, err = template.ParseFiles(
				"web/templates/base.html",
				"web/templates/home/home.html",
				"web/templates/home/_hero.html",
				"web/templates/home/_timeline.html",
				"web/templates/home/_directions.html",
				"web/templates/home/_dresscode.html",
				"web/templates/home/_contact.html",
				"web/templates/home/_rsvp-cta.html",
				"web/templates/home/_rsvp-frame.html",
				"web/templates/home/_sticky-footer.html",
			)
			if err != nil {
				log.Warn("Template reload failed", "err", err)
				http.Error(w, "Template error", http.StatusInternalServerError)
				return
			}
			log.Info("Home templates reloaded (dev mode)")
		}

		// Get invites and pass them to template
		var tplData homeTemplateData
		invites, err := s.ReadAllInvites()
		if err != nil {
			http.Error(w, "Could not get invites", 500)
			return
		}

		tplData.Invites = invites

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")

		if err := homeTemplate.ExecuteTemplate(w, "base", tplData); err != nil {
			log.Error("Template execution failed", "err", err)
		}
	}
}

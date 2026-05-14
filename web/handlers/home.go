package handlers

import (
	"html/template"
	"net/http"

	"github.com/charmbracelet/log"
)

var homeTemplate *template.Template

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

func GetHome(w http.ResponseWriter, r *http.Request) {
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

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	if err := homeTemplate.ExecuteTemplate(w, "base", nil); err != nil {
		log.Error("Template execution failed", "err", err)
	}
}

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
		"public/_home/home.html",
		"public/_home/_hero.html",
		"public/_home/_timeline.html",
		"public/_home/_directions.html",
		"public/_home/_dresscode.html",
		"public/_home/_contact.html",
		"public/_home/_rsvp-cta.html",
		"public/_home/_sticky-footer.html",
	)
	if err != nil {
		log.Debug("Home templates not yet available", "err", err)
	}
}

func getHome(w http.ResponseWriter, r *http.Request) {
	// Reload on error to allow dev editing without restart
	if homeTemplate == nil {
		var err error
		homeTemplate, err = template.ParseFiles(
			"public/_home/home.html",
			"public/_home/_hero.html",
			"public/_home/_timeline.html",
			"public/_home/_directions.html",
			"public/_home/_dresscode.html",
			"public/_home/_contact.html",
			"public/_home/_rsvp-cta.html",
			"public/_home/_sticky-footer.html",
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

	if err := homeTemplate.ExecuteTemplate(w, "home", nil); err != nil {
		log.Error("Template execution failed", "err", err)
	}
}

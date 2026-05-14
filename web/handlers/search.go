package handlers

import (
	"html/template"
	"net/http"

	"github.com/charmbracelet/log"

	"rsvp/invite"
	"rsvp/store"
)

var searchTemplate *template.Template

type searchTemplateData struct {
	Invites []*invite.Invite
}

func init() {
	// Silent fail on init - templates might not exist yet during startup
	var err error
	searchTemplate, err = template.New("search-tpl").Funcs(funcMap).ParseFiles(
		"web/templates/rsvp/search.html",
	)
	if err != nil {
		log.Warn("Home templates not yet available", "err", err)
	}
}

func GetSearch(s *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if searchTemplate == nil {
			var err error

			searchTemplate, err = template.New("search-tpl").Funcs(funcMap).ParseFiles(
				"web/templates/rsvp/search.html",
			)
			if err != nil {
				log.Warn("Template reload failed", "err", err)
				http.Error(w, "Template error", http.StatusInternalServerError)
				return
			}
			log.Info("Search template reloaded (dev mode)")
		}

		query := r.URL.Query().Get("search-input")
		if query == "" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(""))
			return
		}

		log.Info(query)

		// Get invites and pass them to template
		var tplData searchTemplateData
		invites, err := s.SearchInvites(query)
		if err != nil {
			log.Error(err)
			http.Error(w, "Could not get invites", 500)
			return
		}

		tplData.Invites = invites

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")

		if err := searchTemplate.ExecuteTemplate(w, "search-tpl", tplData); err != nil {
			log.Error("Template execution failed", "err", err)
		}
	}
}

package handlers

import (
	"html/template"
	"net/http"
	"slices"

	"github.com/charmbracelet/log"

	"rsvp/config"
	"rsvp/invite"
	"rsvp/store"
)

var overviewTemplate *template.Template

type OverviewStats struct {
	NoResponseDay       int
	NoResponseEvening   int
	AttendingDay        int
	AttendingEvening    int
	NotAttendingDay     int
	NotAttendingEvening int
	InvitedDay          int
	InvitedEvening      int
}

type OverviewData struct {
	Invites []*invite.InviteWithRSVPs
	Stats   OverviewStats
}

var funcMap template.FuncMap = template.FuncMap{
	"last": func(slice any) any {
		switch v := slice.(type) {
		case []*invite.RSVP:
			if len(v) > 0 {
				return v[len(v)-1]
			}
		}
		return nil
	},
}

func init() {
	// Silent fail on init - templates might not exist yet during startup
	var err error
	overviewTemplate, err = template.New("").Funcs(funcMap).ParseFiles(
		"web/templates/base.html",
		"web/templates/overview/overview.html",
		"web/templates/overview/stats.html",
		"web/templates/overview/invite-card.html",
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
			overviewTemplate, err = template.New("").Funcs(funcMap).ParseFiles(
				"web/templates/base.html",
				"web/templates/overview/overview.html",
				"web/templates/overview/stats.html",
				"web/templates/overview/invite-card.html",
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

		// Sort invites alphabetically by name
		slices.SortFunc(data, func(a, b *invite.InviteWithRSVPs) int {
			if a.Invite.Name < b.Invite.Name {
				return -1
			}
			if a.Invite.Name > b.Invite.Name {
				return 1
			}
			return 0
		})

		stats := OverviewStats{}
		// Cycle through data and collect le stats
		for _, invitation := range data {
		
			// Get latest RSVP
			// Do as I say, not as I do, ok? Using this nil check is GROSS and
			// not idiomatic Go but it avoids doing the sort twice.
			var latest *invite.RSVP = nil
			if len(invitation.RSVPs) > 0 {
				latest = slices.MaxFunc(invitation.RSVPs, func(a, b *invite.RSVP) int {
					return a.Timestamp.Compare(b.Timestamp)
				})
			}

			// Handle day
			if invitation.Invite.Day {
				stats.InvitedDay++

				if latest == nil {
					stats.NoResponseDay++
				} else if latest.AttendingDay {
					stats.AttendingDay++
				} else {
					stats.NotAttendingDay++
				}
			}

			// All invites are for the evening by default
			stats.InvitedEvening++

			if latest == nil {
				stats.NoResponseEvening++
			} else if latest.AttendingEvening {
				stats.AttendingEvening++
			} else {
				stats.NotAttendingEvening++
			}
		}

		templateData := OverviewData{
			Invites: data,
			Stats:   stats,
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

package handlers

import (
	"html/template"
	"net/http"

	"github.com/charmbracelet/log"
	"rsvp/domain"
)

// InvitesPageData holds the data to be rendered on the invites page
type InvitesPageData struct {
	Invites          []*domain.InviteWithRSVPs
	DayAttendees     int
	EveningAttendees int
}

var invitesTemplate *template.Template

func init() {
	// Silent fail on init - templates might not exist yet during startup
	var err error
	invitesTemplate, err = template.ParseFiles("public/_invites/invites.html")
	if err != nil {
		log.Debug("Invites templates not yet available", "err", err)
	}
}

func getInvites(w http.ResponseWriter, r *http.Request) {
	log.Info("GET /invite")

	s, err := StoreFromContext(r.Context())
	if err != nil {
		log.Error("store not available", "err", err)
		http.Error(w, "Store unavailable", http.StatusInternalServerError)
		return
	}

	// Fetch all invites with their RSVPs from the store
	invites, err := s.ReadAllInvitesWithRSVPs()
	if err != nil {
		log.Error("failed to read invites", "err", err)
		http.Error(w, "Failed to load invites", http.StatusInternalServerError)
		return
	}

	// Calculate summary statistics based on most recent RSVP per invite
	dayAttendees := 0
	eveningAttendees := 0
	for _, inviteWithRsvps := range invites {
		latestRsvp := inviteWithRsvps.LatestRSVP()
		if latestRsvp != nil {
			for _, attendance := range latestRsvp.Attendances {
				if attendance.AttendingDay {
					dayAttendees++
				}
				if attendance.AttendingEvening {
					eveningAttendees++
				}
			}
		}
	}

	// Prepare page data with invites and summary stats
	pageData := InvitesPageData{
		Invites:          invites,
		DayAttendees:     dayAttendees,
		EveningAttendees: eveningAttendees,
	}

	// Reload on error to allow dev editing without restart
	if invitesTemplate == nil {
		var err error
		invitesTemplate, err = template.ParseFiles("public/_invites/invites.html")
		if err != nil {
			log.Warn("Template reload failed", "err", err)
			http.Error(w, "Template error", http.StatusInternalServerError)
			return
		}
		log.Info("Invites templates reloaded (dev mode)")
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := invitesTemplate.Execute(w, pageData); err != nil {
		log.Error("Template execution failed", "err", err)
	}
}

func getInvite(w http.ResponseWriter, r *http.Request) {
	inviteId := r.PathValue("inviteId")
	log.Info("GET /invite/{inviteId}", "inviteId", inviteId)

	s, err := StoreFromContext(r.Context())
	if err != nil {
		log.Error("store not available", "err", err)
		http.Error(w, "Store unavailable", http.StatusInternalServerError)
		return
	}

	_ = s // TODO: Use store to fetch invite
	w.Write([]byte("TODO: Render invite template"))
}

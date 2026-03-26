package handlers

import (
	"html/template"
	"net/http"

	"github.com/charmbracelet/log"
	"rsvp/domain"
)

var respondTemplate *template.Template
var respondFormTemplate *template.Template
var respondSuccessTemplate *template.Template

func init() {
	// Silent fail on init - templates might not exist yet during startup
	var err error
	respondTemplate, err = template.ParseFiles("public/_respond/respond.html")
	if err != nil {
		log.Debug("Respond templates not yet available", "err", err)
	}

	respondFormTemplate, err = template.ParseFiles("public/_respond/respond-form-fields.html")
	if err != nil {
		log.Debug("Respond form templates not yet available", "err", err)
	}

	respondSuccessTemplate, err = template.ParseFiles("public/_success/respond-success.html")
	if err != nil {
		log.Debug("Success templates not yet available", "err", err)
	}
}

func getRespond(w http.ResponseWriter, r *http.Request) {
	log.Info("GET /respond")

	s, err := StoreFromContext(r.Context())
	if err != nil {
		log.Error("store not available", "err", err)
		http.Error(w, "Store unavailable", http.StatusInternalServerError)
		return
	}

	// Fetch all invites for the dropdown
	invites, err := s.ReadAllInvitesWithRSVPs()
	if err != nil {
		log.Error("failed to read invites", "err", err)
		http.Error(w, "Failed to load invites", http.StatusInternalServerError)
		return
	}

	// Reload on error to allow dev editing without restart
	if respondTemplate == nil {
		var err error
		respondTemplate, err = template.ParseFiles("public/_respond/respond.html")
		if err != nil {
			log.Warn("Template reload failed", "err", err)
			http.Error(w, "Template error", http.StatusInternalServerError)
			return
		}
		log.Info("Respond templates reloaded (dev mode)")
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := respondTemplate.Execute(w, invites); err != nil {
		log.Error("Template execution failed", "err", err)
	}
}

func getRespondFormFields(w http.ResponseWriter, r *http.Request) {
	inviteId := r.URL.Query().Get("inviteId")
	log.Info("GET /respond/form", "inviteId", inviteId)

	s, err := StoreFromContext(r.Context())
	if err != nil {
		log.Error("store not available", "err", err)
		http.Error(w, "Store unavailable", http.StatusInternalServerError)
		return
	}

	// Fetch the specific invite
	inviteWithRsvps, err := s.ReadInviteWithRSVPs(inviteId)
	if err != nil {
		log.Error("failed to read invite", "inviteId", inviteId, "err", err)
		http.Error(w, "Failed to load invite", http.StatusInternalServerError)
		return
	}

	// Reload on error to allow dev editing without restart
	if respondFormTemplate == nil {
		var err error
		respondFormTemplate, err = template.ParseFiles("public/_respond/respond-form-fields.html")
		if err != nil {
			log.Warn("Template reload failed", "err", err)
			http.Error(w, "Template error", http.StatusInternalServerError)
			return
		}
		log.Info("Respond form templates reloaded (dev mode)")
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := respondFormTemplate.Execute(w, inviteWithRsvps); err != nil {
		log.Error("Template execution failed", "err", err)
	}
}

func postRsvp(w http.ResponseWriter, r *http.Request) {
	inviteId := r.PathValue("inviteId")
	log.Info("POST /rsvp", "id", inviteId)

	s, err := StoreFromContext(r.Context())
	if err != nil {
		log.Error("store not available", "err", err)
		http.Error(w, "Store unavailable", http.StatusInternalServerError)
		return
	}

	// Parse form data
	err = r.ParseForm()
	if err != nil {
		log.Error("could not parse form", "err", err)
		http.Error(w, "Could not parse form", http.StatusBadRequest)
		return
	}

	// Create RSVP from form data
	rsvp, err := domain.NewRsvp(inviteId)
	if err != nil {
		log.Error("could not create rsvp for inviteId", "inviteId", inviteId, "err", err)
		http.Error(w, "Invalid invite ID", http.StatusBadRequest)
		return
	}

	// Parse form values and build Attendances array
	rsvp.Attendances = []domain.PersonAttendance{
		{
			PersonIndex:      0,
			AttendingDay:     r.FormValue("person0AttendingDay") == "on",
			AttendingEvening: r.FormValue("person0AttendingEvening") == "on",
		},
	}

	// Check if this invite has a second person
	if r.FormValue("hasPerson1") == "true" {
		rsvp.Attendances = append(rsvp.Attendances, domain.PersonAttendance{
			PersonIndex:      1,
			AttendingDay:     r.FormValue("person1AttendingDay") == "on",
			AttendingEvening: r.FormValue("person1AttendingEvening") == "on",
		})
	}

	rsvp.HasPresentation = r.FormValue("hasPresentation") == "on"
	rsvp.PhoneNumber = r.FormValue("phoneNumber")
	rsvp.DietNotes = r.FormValue("dietNotes")
	rsvp.Message = r.FormValue("message")

	log.Info("RSVP submitted", "invite", inviteId, "rsvp", rsvp)

	// TODO: Save RSVP to database
	_ = s

	// Fetch the invite and return success page
	inviteWithRsvps, err := s.ReadInviteWithRSVPs(inviteId)
	if err != nil {
		log.Error("failed to read invite", "inviteId", inviteId, "err", err)
		http.Error(w, "Failed to load invite", http.StatusInternalServerError)
		return
	}

	// Data for success template
	successData := map[string]interface{}{
		"Invite": inviteWithRsvps.Invite,
		"RSVP":   rsvp,
	}

	// Reload on error to allow dev editing without restart
	if respondSuccessTemplate == nil {
		var err error
		respondSuccessTemplate, err = template.ParseFiles("public/_success/respond-success.html")
		if err != nil {
			log.Warn("Template reload failed", "err", err)
			http.Error(w, "Template error", http.StatusInternalServerError)
			return
		}
		log.Info("Success templates reloaded (dev mode)")
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := respondSuccessTemplate.Execute(w, successData); err != nil {
		log.Error("Template execution failed", "err", err)
	}
}

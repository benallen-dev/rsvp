package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Create/open test database
	db, err := sql.Open("sqlite3", "test-data.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create tables
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS invites (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			day BOOLEAN,
			evening BOOLEAN
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS rsvps (
			id TEXT PRIMARY KEY,
			invite_id TEXT NOT NULL,
			timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
			attending_day BOOLEAN,
			attending_evening BOOLEAN,
			has_presentation BOOLEAN,
			phone_number TEXT,
			diet_notes TEXT,
			message TEXT,
			FOREIGN KEY (invite_id) REFERENCES invites(id)
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	// Clear existing data
	db.Exec("DELETE FROM rsvps")
	db.Exec("DELETE FROM invites")

	// Create invites
	invites := []struct {
		name    string
		day     bool
		evening bool
	}{
		{"Sarah Johnson (Evening)", false, true},
		{"Michael Chen", true, true},
		{"Emma Williams", true, true},
		{"James Rodriguez (Evening)", false, true},
		{"Lisa Anderson", true, true},
	}

	inviteIDs := make([]string, len(invites))

	for i, inv := range invites {
		id := uuid.New().String()
		inviteIDs[i] = id

		_, err = db.Exec(
			"INSERT INTO invites (id, name, day, evening) VALUES (?, ?, ?, ?)",
			id, inv.name, inv.day, inv.evening,
		)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Created invite: %s (%s)\n", inv.name, id)
	}

	// Create RSVPs with realistic distribution
	rsvpData := []struct {
		inviteIdx        int
		attendingDay     bool
		attendingEvening bool
		hasPresentation  bool
		phoneNumber      string
		dietNotes        string
		message          string
		hoursOffset      int
	}{
		// Invite 0 (Sarah): 0 RSVPs
		// Invite 1 (Michael): 1 RSVP
		{1, true, true, false, "", "Vegetarian", "Really looking forward to it!", 24},
		// Invite 2 (Emma): 3 RSVPs
		{2, true, false, true, "555-0123", "Gluten-free", "See you then!", 48},
		{2, false, true, false, "", "", "Can't wait for the evening!", 72},
		{2, true, true, true, "555-0456", "No shellfish", "", 120},
		// Invite 3 (James): 2 RSVPs
		{3, true, false, false, "", "Vegan", "Thanks for the invite!", 36},
		{3, true, false, true, "555-0789", "", "Looking forward to it", 60},
		// Invite 4 (Lisa): 1 RSVP
		{4, true, true, false, "", "", "Excited to attend!", 96},
	}

	baseTime := time.Now().AddDate(0, 0, -7) // 7 days ago

	for i, rsvp := range rsvpData {
		rsvpID := uuid.New().String()
		timestamp := baseTime.Add(time.Duration(rsvp.hoursOffset) * time.Hour)

		_, err = db.Exec(
			"INSERT INTO rsvps (id, invite_id, timestamp, attending_day, attending_evening, has_presentation, phone_number, diet_notes, message) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
			rsvpID,
			inviteIDs[rsvp.inviteIdx],
			timestamp.Format(time.RFC3339),
			rsvp.attendingDay,
			rsvp.attendingEvening,
			rsvp.hasPresentation,
			rsvp.phoneNumber,
			rsvp.dietNotes,
			rsvp.message,
		)
		if err != nil {
			log.Fatal(err)
		}
		presentationStr := ""
		if rsvp.hasPresentation {
			presentationStr = " with presentation"
		}
		fmt.Printf("Created RSVP %d: %s attending (day: %v, evening: %v)%s\n",
			i+1, invites[rsvp.inviteIdx].name, rsvp.attendingDay, rsvp.attendingEvening, presentationStr)
	}

	fmt.Println("\n✓ Test database created successfully at test-data.db")
	fmt.Println("Summary:")
	fmt.Println("  - 5 invites")
	fmt.Println("  - 7 RSVPs total")
	fmt.Println("  - Distribution: 0, 1, 3, 2, 1 RSVPs per invite")
}

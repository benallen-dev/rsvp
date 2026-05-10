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

	// TODO: Use actual structs from invite package to at least somewhat keep this in sync
	// Create invites
	invites := []struct {
		name    string
		day     bool
		evening bool
	}{
		{"Sarah Johnson", true, true},
		{"Michael Chen", true, true},
		{"Emma Williams", true, true},
		{"James Rodriguez", false, true},
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
		dietNotes        string
		message          string
		hoursOffset      int
	}{
		// Invite 0 (Sarah): 0 RSVPs
		// Invite 1 (Michael): 1 RSVP
		{1, true, true, "Vegetarian", "Really looking forward to it!", 24},
		// Invite 2 (Emma): 3 RSVPs
		{2, true, false, "Gluten-free", "See you then!", 48},
		{2, false, true, "", "Can't wait for the evening!", 72},
		{2, true, true, "No shellfish", "", 120},
		// Invite 3 (James): 2 RSVPs
		{3, true, false, "Vegan", "Thanks for the invite!", 36},
		{3, true, false, "", "Looking forward to it", 60},
		// Invite 4 (Lisa): 1 RSVP
		{4, true, true, "", "Excited to attend!", 96},
	}

	baseTime := time.Now().AddDate(0, 0, -7) // 7 days ago

	for i, rsvp := range rsvpData {
		rsvpID := uuid.New().String()
		timestamp := baseTime.Add(time.Duration(rsvp.hoursOffset) * time.Hour)

		_, err = db.Exec(
			"INSERT INTO rsvps (id, invite_id, timestamp, attending_day, attending_evening, diet_notes, message) VALUES (?, ?, ?, ?, ?, ?, ?)",
			rsvpID,
			inviteIDs[rsvp.inviteIdx],
			timestamp.Format(time.RFC3339),
			rsvp.attendingDay,
			rsvp.attendingEvening,
			rsvp.dietNotes,
			rsvp.message,
		)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Created RSVP %d: %s attending (day: %v, evening: %v)\n",
			i+1, invites[rsvp.inviteIdx].name, rsvp.attendingDay, rsvp.attendingEvening)
	}

	fmt.Println("\n✓ Test database created successfully at test-data.db")
	fmt.Println("Summary:")
	fmt.Println("  - 5 invites")
	fmt.Println("  - 7 RSVPs total")
	fmt.Println("  - Distribution: 0, 1, 3, 2, 1 RSVPs per invite")
}

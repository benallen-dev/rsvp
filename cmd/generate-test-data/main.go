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

	// Create tables with new schema
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS invites (
			id TEXT PRIMARY KEY,
			person_0_first_name TEXT NOT NULL,
			person_0_last_name TEXT NOT NULL,
			person_1_first_name TEXT,
			person_1_last_name TEXT,
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
			person_0_attending_day BOOLEAN,
			person_0_attending_evening BOOLEAN,
			person_1_attending_day BOOLEAN,
			person_1_attending_evening BOOLEAN,
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
		people []struct {
			firstName string
			lastName  string
		}
		day     bool
		evening bool
	}{
		{
			people: []struct {
				firstName string
				lastName  string
			}{
				{firstName: "Sarah", lastName: "Johnson"},
				{firstName: "Michael", lastName: "Johnson"},
			},
			day: false, evening: true,
		},
		{
			people: []struct {
				firstName string
				lastName  string
			}{
				{firstName: "David", lastName: "Chen"},
				{firstName: "Jennifer", lastName: "Chen"},
			},
			day: true, evening: true,
		},
		{
			people: []struct {
				firstName string
				lastName  string
			}{
				{firstName: "Emma", lastName: "Williams"},
			},
			day: true, evening: true,
		},
		{
			people: []struct {
				firstName string
				lastName  string
			}{
				{firstName: "James", lastName: "Rodriguez"},
				{firstName: "Maria", lastName: "Garcia"},
			},
			day: false, evening: true,
		},
		{
			people: []struct {
				firstName string
				lastName  string
			}{
				{firstName: "Lisa", lastName: "Anderson"},
				{firstName: "Tom", lastName: "Anderson"},
			},
			day: true, evening: true,
		},
	}

	inviteIDs := make([]string, len(invites))

	for i, inv := range invites {
		id := uuid.New().String()
		inviteIDs[i] = id

		person0FirstName := inv.people[0].firstName
		person0LastName := inv.people[0].lastName
		var person1FirstName, person1LastName *string
		if len(inv.people) > 1 {
			person1FirstName = &inv.people[1].firstName
			person1LastName = &inv.people[1].lastName
		}

		_, err = db.Exec(
			"INSERT INTO invites (id, person_0_first_name, person_0_last_name, person_1_first_name, person_1_last_name, day, evening) VALUES (?, ?, ?, ?, ?, ?, ?)",
			id, person0FirstName, person0LastName, person1FirstName, person1LastName, inv.day, inv.evening,
		)
		if err != nil {
			log.Fatal(err)
		}

		// Display name formatting
		displayName := person0FirstName + " " + person0LastName
		if len(inv.people) > 1 {
			if inv.people[1].lastName == person0LastName {
				displayName = person0FirstName + " & " + inv.people[1].firstName + " " + person0LastName
			} else {
				displayName = person0FirstName + " " + person0LastName + " & " + inv.people[1].firstName + " " + inv.people[1].lastName
			}
		}
		fmt.Printf("Created invite: %s (%s)\n", displayName, id)
	}

	// Create RSVPs with realistic distribution
	rsvpData := []struct {
		inviteIdx               int
		person0AttendingDay     bool
		person0AttendingEvening bool
		person1AttendingDay     bool
		person1AttendingEvening bool
		hasPresentation         bool
		phoneNumber             string
		dietNotes               string
		message                 string
		hoursOffset             int
	}{
		// Invite 0 (Sarah & Michael): 0 RSVPs
		// Invite 1 (David & Jennifer): 1 RSVP
		{1, true, true, true, true, false, "", "Vegetarian", "Really looking forward to it!", 24},
		// Invite 2 (Emma): 3 RSVPs
		{2, true, false, false, false, true, "555-0123", "Gluten-free", "See you then!", 48},
		{2, false, true, false, false, false, "", "", "Can't wait for the evening!", 72},
		{2, true, true, false, false, true, "555-0456", "No shellfish", "", 120},
		// Invite 3 (James & Maria): 2 RSVPs
		{3, true, false, false, true, false, "", "Vegan", "Thanks for the invite!", 36},
		{3, true, false, true, false, true, "555-0789", "", "Looking forward to it", 60},
		// Invite 4 (Lisa & Tom): 1 RSVP
		{4, true, true, false, true, false, "", "", "Excited to attend!", 96},
	}

	baseTime := time.Now().AddDate(0, 0, -7) // 7 days ago

	for i, rsvp := range rsvpData {
		rsvpID := uuid.New().String()
		timestamp := baseTime.Add(time.Duration(rsvp.hoursOffset) * time.Hour)

		_, err = db.Exec(
			"INSERT INTO rsvps (id, invite_id, timestamp, person_0_attending_day, person_0_attending_evening, person_1_attending_day, person_1_attending_evening, has_presentation, phone_number, diet_notes, message) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
			rsvpID,
			inviteIDs[rsvp.inviteIdx],
			timestamp.Format(time.RFC3339),
			rsvp.person0AttendingDay,
			rsvp.person0AttendingEvening,
			rsvp.person1AttendingDay,
			rsvp.person1AttendingEvening,
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
		fmt.Printf("Created RSVP %d: invite %d (person0 - day: %v, evening: %v; person1 - day: %v, evening: %v)%s\n",
			i+1, rsvp.inviteIdx, rsvp.person0AttendingDay, rsvp.person0AttendingEvening, rsvp.person1AttendingDay, rsvp.person1AttendingEvening, presentationStr)
	}

	fmt.Println("\n✓ Test database created successfully at test-data.db")
	fmt.Println("Summary:")
	fmt.Println("  - 5 invites (4 pairs + 1 single)")
	fmt.Println("  - 7 RSVPs total")
	fmt.Println("  - Distribution: 0, 1, 3, 2, 1 RSVPs per invite")
}

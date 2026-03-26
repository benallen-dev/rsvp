package main

import (
	"database/sql"
	"fmt"
	"log"

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
		firstName string
		lastName  string
		day       bool
		evening   bool
	}{
		// Day guests
		{"Karin", "Gronsveld", true, true},
		{"Marijke", "Valkenburg", true, true},
		{"Simon", "Knot", true, true},
		{"Diana", "Hendriks", true, true},
		{"Rod", "Allen", true, true},
		{"Marja", "Allen", true, true},
		{"Alexander", "Allen", true, true},
		{"Romy", "Schoemakers", true, true},
		{"Victoria", "Moek", true, true},
		{"Kevin", "Bouhuizen", true, true},
		{"Nina", "van Huizen", true, true},
		{"Charley", "Bakker", true, true},
		{"Bram", "van Huizen", true, true},
		{"Denise", "Komen", true, true},
		{"Rafael", "Endeman", true, true},
		{"Joost", "Kerpels", true, true},
		{"Naomi", "du Pree", true, true},
		{"Jeroen", "Ouweneel", true, true},
		{"Matthijs", "Alderliesten", true, true},
		{"Roel", "de Rijk", true, true},
		{"Matthijs", "Weskin", true, true},
		{"Reinier", "Zwikker", true, true},
		{"Tom", "van Dijk", true, true},
		{"Tim", "Velzeboer", true, true},
		{"Chris", "Bootsman", true, true},
		{"Tim", "Verzeide", true, true},
		{"Natasja", "Oostrum", true, true},
		// Evening guests
		{"Jolien", "Steenweg", false, true},
		{"Sabine", "de Richemont", false, true},
		{"Hans", "Gouw", false, true},
		{"Hilda", "Gouw", false, true},
		{"Jan-Willem", "'t Hart", false, true},
		{"Bonnie", "'t Hart", false, true},
		{"Henriëtte", "Knot", false, true},
		{"Joey", "Kempff", false, true},
		{"Vera", "Gouw", false, true},
		{"Jasper", "Hentzen", false, true},
		{"Anita", "van der Ploeg", false, true},
		{"Merel", "Demenint", false, true},
		{"Niels", "Heisterkamp", false, true},
		{"Michelle", "(Niels +1)", false, true},
		{"Mo", "de Ruijter", false, true},
		{"Nieneke", "de Ruijter", false, true},
		{"Job", "van der Weide", false, true},
		{"Shelley", "Rijnenberg", false, true},
		{"Melany", "Thunnissen", false, true},
		{"Keyla", "Cornes-Maduro", false, true},
		{"Joël", "Kema", false, true},
		{"Nicholas", "Molenaar", false, true},
		{"Sander", "Liebens", false, true},
		{"Nynke", "de Rover", false, true},
		{"Erik", "Gronsveld", false, true},
		{"Monica", "Gronsveld", false, true},
		{"Marit", "Gronsveld", false, true},
		{"Joseph", "Verburg", false, true},
		{"Jennie", "Christiaanse", false, true},
		{"Mandy", "Klapwijk", false, true},
		{"Menno", "van Lopik", false, true},
		{"Dennis", "Roepman", false, true},
		{"Cindy", "Roepman", false, true},
		{"Silvano", "Roepman", false, true},
		{"Tony", "Roepman", false, true},
		{"Ilse", "Roepman", false, true},
		{"Kevin", "Roepman", false, true},
		{"Sabine", "Roepman", false, true},
		{"Britt", "Hendriks", false, true},
		{"Jenny", "Taube", false, true},
		{"Alison", "Taube", false, true},
		{"Kasper", "van Steveninck", false, true},
		{"Kasper", "+ 1", false, true},
		{"Michel", "", false, true},
		{"Tessy", "", false, true},
		{"Jasper", "Boot", false, true},
		{"Wendy", "", false, true},
		{"Rico", "", false, true},
		{"Mea", "", false, true},
		{"Sjoerd", "", false, true},
		{"Maartje", "", false, true},
		{"Cynthia", "Slotboom", false, true},
		{"Marco", "Gronsveld", false, true},
		{"Wouter", "van Mierlo", false, true},
		{"Manon", "Speulman", false, true},
		{"Nienke", "van Dam", false, true},
		{"Elise", "Hoffman", false, true},
		{"Jitske", "", false, true},
		{"Kiki", "", false, true},
		{"Cath", "", false, true},
		{"Lianne", "", false, true},
		{"Pepijn", "", false, true},
		{"Daniel", "Melton", false, true},
		{"Mia", "Urem", false, true},
	}

	inviteIDs := make([]string, len(invites))

	for i, inv := range invites {
		id := uuid.New().String()
		inviteIDs[i] = id

		_, err = db.Exec(
			"INSERT INTO invites (id, person_0_first_name, person_0_last_name, person_1_first_name, person_1_last_name, day, evening) VALUES (?, ?, ?, ?, ?, ?, ?)",
			id, inv.firstName, inv.lastName, nil, nil, inv.day, inv.evening,
		)
		if err != nil {
			log.Fatal(err)
		}
		displayName := inv.firstName + " " + inv.lastName
		fmt.Printf("Created invite: %s (%s)\n", displayName, id)
	}

	fmt.Println("\n✓ Test database created successfully at test-data.db")
	fmt.Println("Summary:")
	fmt.Println("  - 86 total invites")
	fmt.Println("  - 27 day guests (day + evening)")
	fmt.Println("  - 59 evening guests (evening only)")
}

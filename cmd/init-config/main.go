package main

import (
	"github.com/charmbracelet/log"

	"rsvp/config"
)

func main() {
	log.Info("Initialising config file")

	// Will crash because of the missing config file on init right now

	newConfig := config.Config{
		LogLevel:    "info",
		AuthEnabled: false,
		Users: config.Users{
			Admin: config.User{
				PrettyName: "Admin",
				Username:   "admin",
				Password:   "admin",
			},
			SuperAdmin: config.User{
				PrettyName: "Super Admin",
				Username:   "superadmin",
				Password:   "superadmin",
			},
		},
	}

	config.Write(newConfig)

	log.Info("Wrote new config to disk")
}

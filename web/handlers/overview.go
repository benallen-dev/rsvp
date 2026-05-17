package handlers

import (
	"net/http"

	"github.com/charmbracelet/log"

	"rsvp/config"
	// "rsvp/invite"
	"rsvp/store"
)

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

		rsvps, err := s.ReadAllRSVPs()
		if err != nil {
			log.Error("could not read rsvps", "err", err)
			http.Error(w, "Data error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

		html := `<!DOCTYPE html>
<html>
<head>
	<meta charset="utf-8">
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<title>RSVP Overview</title>
	<style>
		body { font-family: sans-serif; margin: 20px; }
		table { border-collapse: collapse; width: 100%; }
		th, td { border: 1px solid #ddd; padding: 12px; text-align: left; }
		th { background-color: #f0f0f0; font-weight: bold; }
		tr:nth-child(even) { background-color: #f9f9f9; }
		.true { color: green; }
		.false { color: red; }
	</style>
</head>
<body>
	<h1>RSVP Overview</h1>
	<p>Total RSVPs: <strong>` + string(rune(len(rsvps))) + `</strong></p>
	<table>
		<thead>
			<tr>
				<th>Name</th>
				<th>Type</th>
				<th>Attending Day</th>
				<th>Attending Evening</th>
				<th>Diet Notes</th>
				<th>Message</th>
				<th>Timestamp</th>
			</tr>
		</thead>
		<tbody>`

		for _, rsvp := range rsvps {
			dayClass := "false"
			if rsvp.AttendingDay {
				dayClass = "true"
			}
			eveningClass := "false"
			if rsvp.AttendingEvening {
				eveningClass = "true"
			}

			html += `<tr>
				<td>` + rsvp.Name + `</td>
				<td>` + rsvp.RsvpType + `</td>
				<td class="` + dayClass + `">` + boolStr(rsvp.AttendingDay) + `</td>
				<td class="` + eveningClass + `">` + boolStr(rsvp.AttendingEvening) + `</td>
				<td>` + rsvp.DietNotes + `</td>
				<td>` + rsvp.Message + `</td>
				<td>` + rsvp.Timestamp.Format("2006-01-02 15:04:05") + `</td>
			</tr>`
		}

		html += `</tbody>
	</table>
</body>
</html>`

		w.Write([]byte(html))
	}
}

func boolStr(b bool) string {
	if b {
		return "Yes"
	}
	return "No"
}

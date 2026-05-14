package handlers

import (
	"html/template"

	"rsvp/invite"
)

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
	"sub": func(a, b int) int { return a - b },
}

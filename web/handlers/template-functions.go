package handlers

import (
	"fmt"
	"html/template"
	"time"

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
	"dict": func(values ...any) (map[string]any, error) {
		if len(values)%2 != 0 {
			return nil, fmt.Errorf("dict requires even number of arguments")
		}
		m := make(map[string]any)
		for i := 0; i < len(values); i += 2 {
			key := values[i].(string)
			m[key] = values[i+1]
		}
		return m, nil
	},
	"formatDate": func(t time.Time) string {
		months := []string{"januari", "februari", "maart", "april", "mei", "juni", "juli", "augustus", "september", "oktober", "november", "december"}

		return fmt.Sprintf("%d %s %d", t.Day(), months[t.Month()-1], t.Year())
	},
}

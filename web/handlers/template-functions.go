package handlers

import (
	"fmt"
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
}

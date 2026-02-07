package main

import (
	"fmt"
	"net/http"
	"rsvp/web"
)

const PORT = ":8080"

func main() {
	mux := web.NewMux()

	fmt.Println("listening on " + PORT)
	http.ListenAndServe(PORT, mux)
}

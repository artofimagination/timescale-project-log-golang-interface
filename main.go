package main

import (
	"net/http"

	"github.com/artofimagination/timescaledb-project-log-go-interface/restcontrollers"
)

func main() {
	_, err := restcontrollers.NewRESTController()
	if err != nil {
		panic(err)
	}

	// Start HTTP server that accepts requests from the offer process to exchange SDP and Candidates
	panic(http.ListenAndServe(":8186", nil))
}

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/artofimagination/timescaledb-project-log-go-interface/dbcontrollers"
	"github.com/artofimagination/timescaledb-project-log-go-interface/timescaledb"

	"github.com/pkg/errors"
)

func sayHello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hi! I am an example server!")
}

const (
	POST = "POST"
	GET  = "GET"
)

func checkRequestType(requestTypeString string, w http.ResponseWriter, r *http.Request) error {
	if r.Method != requestTypeString {
		w.WriteHeader(http.StatusBadRequest)
		errorString := fmt.Sprintf("Invalid request type %s", r.Method)
		fmt.Fprint(w, errorString)
		return errors.New(errorString)
	}
	return nil
}

func decodePostData(w http.ResponseWriter, r *http.Request) (map[string]interface{}, error) {
	if err := checkRequestType(POST, w, r); err != nil {
		return nil, err
	}

	data := make(map[string]interface{})
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err = errors.Wrap(errors.WithStack(err), "Failed to decode request json")
		fmt.Fprint(w, err.Error())
		return nil, err
	}

	return data, nil
}

func addProject(w http.ResponseWriter, r *http.Request) {
	log.Println("Adding project")
	_, err := decodePostData(w, r)
	if err != nil {
		return
	}

}

func getProject(w http.ResponseWriter, r *http.Request) {
	log.Println("Getting project")
	if err := checkRequestType(GET, w, r); err != nil {
		return
	}
}

func main() {
	http.HandleFunc("/", sayHello)
	http.HandleFunc("/add-project", addProject)
	http.HandleFunc("/get-project", getProject)

	_, err := dbcontrollers.NewDBController()
	if err != nil {
		panic(err)
	}

	if err := timescaledb.BootstrapData(); err != nil {
		log.Fatalf("Data bootstrap failed. %s\n", errors.WithStack(err))
	}

	// Start HTTP server that accepts requests from the offer process to exchange SDP and Candidates
	panic(http.ListenAndServe(":8080", nil))
}

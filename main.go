package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/artofimagination/timescaledb-project-log-go-interface/dbcontrollers"
	"github.com/artofimagination/timescaledb-project-log-go-interface/models"

	"github.com/pkg/errors"
)

var controller *dbcontrollers.TimescaleController

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

func addData(w http.ResponseWriter, r *http.Request) {
	log.Println("Adding data")
	input, err := decodePostData(w, r)
	if err != nil {
		return
	}

	inputDataList, ok := input["data_to_store"].([]interface{})
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Missing 'data_to_store'")
		return
	}

	dataList := make([]models.Data, len(inputDataList))
	for i, data := range inputDataList {
		viewerID, ok := data.(map[string]interface{})["viewer_id"].(int)
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Missing 'viewer_id'")
			return
		}

		content, ok := data.(map[string]interface{})["data"].(models.DataMap)
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Missing 'data'")
			return
		}
		d := models.Data{
			ViewerID: viewerID,
			Data:     content,
		}
		dataList[i] = d
	}

	if err := controller.AddData(dataList); err == nil {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Data addition complete")
		return
	}
	if err.Error() == dbcontrollers.ErrFailedToAddData.Error() {
		w.WriteHeader(http.StatusAccepted)
		fmt.Fprint(w, err.Error())
		return
	}
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprint(w, err.Error())
}

func getDataByViewer(w http.ResponseWriter, r *http.Request) {
	log.Println("Getting data by viewer")
	if err := checkRequestType(GET, w, r); err != nil {
		return
	}
}

func main() {
	http.HandleFunc("/", sayHello)
	http.HandleFunc("/add-data", addData)
	http.HandleFunc("/get-data-by-viewer", getDataByViewer)

	c, err := dbcontrollers.NewDBController()
	if err != nil {
		panic(err)
	}
	controller = c

	// Start HTTP server that accepts requests from the offer process to exchange SDP and Candidates
	panic(http.ListenAndServe(":8080", nil))
}

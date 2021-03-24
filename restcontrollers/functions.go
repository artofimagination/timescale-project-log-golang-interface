package restcontrollers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/artofimagination/timescaledb-project-log-go-interface/dbcontrollers"
	"github.com/pkg/errors"
)

type RESTController struct {
	DBController *dbcontrollers.TimescaleController
}

const (
	DataPathAdd             = "/add-data"
	DataPathGetDataByViewer = "/get-data-by-viewer"
)

const (
	POST = "POST"
	GET  = "GET"
)

var DataOK = "\"OK\""

type ResponseWriter struct {
	http.ResponseWriter
}

type Request struct {
	*http.Request
}

type ResponseData struct {
	Error string      `json:"error" validation:"required"`
	Data  interface{} `json:"data" validation:"required"`
}

func (w ResponseWriter) writeError(message string, statusCode int) {
	response := ResponseData{
		Error: message,
	}

	b, err := json.Marshal(response)
	if err != nil {
		w.writeError(err.Error(), http.StatusInternalServerError)
		return
	}
	w.writeResponse(string(b), statusCode)
}

func (w ResponseWriter) writeData(data interface{}, statusCode int) {
	response := ResponseData{
		Data: data,
	}

	b, err := json.Marshal(response)
	if err != nil {
		w.writeError(err.Error(), http.StatusInternalServerError)
		return
	}
	w.writeResponse(string(b), statusCode)
}

func (w ResponseWriter) writeResponse(data string, statusCode int) {
	w.WriteHeader(statusCode)
	fmt.Fprint(w, data)
}

func checkRequestType(requestTypeString string, w ResponseWriter, r *Request) error {
	if r.Method != requestTypeString {
		w.WriteHeader(http.StatusBadRequest)
		errorString := fmt.Sprintf("Invalid request type %s", r.Method)
		return errors.New(errorString)
	}
	return nil
}

func decodePostData(w ResponseWriter, r *Request) (map[string]interface{}, error) {
	if err := checkRequestType(POST, w, r); err != nil {
		return nil, err
	}

	data := make(map[string]interface{})
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err = errors.Wrap(errors.WithStack(err), "Failed to decode request json")
		return nil, err
	}

	return data, nil
}

func sayHello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hi! I am a project log database server!")
}

func makeHandler(fn func(ResponseWriter, *Request)) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		r := &Request{request}
		w := ResponseWriter{writer}
		fn(w, r)
	}
}

func NewRESTController() (*RESTController, error) {
	dbController, err := dbcontrollers.NewDBController()
	if err != nil {
		return nil, err
	}

	restController := &RESTController{
		DBController: dbController,
	}

	http.HandleFunc("/", sayHello)
	http.HandleFunc(DataPathAdd, makeHandler(restController.addData))
	http.HandleFunc(DataPathGetDataByViewer, makeHandler(restController.getDataByViewer))

	return restController, nil
}

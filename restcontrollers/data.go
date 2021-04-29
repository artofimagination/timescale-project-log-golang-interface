package restcontrollers

import (
	"log"
	"net/http"
	"time"

	"github.com/artofimagination/timescaledb-project-log-go-interface/dbcontrollers"
	"github.com/artofimagination/timescaledb-project-log-go-interface/models"
	"github.com/google/uuid"
)

func (c *RESTController) addData(w ResponseWriter, r *Request) {
	log.Println("Adding data")
	input, err := decodePostData(w, r)
	if err != nil {
		return
	}

	inputDataList, ok := input["data_to_store"].([]interface{})
	if !ok {
		w.writeError("Missing 'data_to_store'", http.StatusBadRequest)
		return
	}

	dataList := make([]models.Data, len(inputDataList))
	for i, data := range inputDataList {
		layout := "Mon Jan 02 2006 15:04:05.0000 GMT-0700"
		createAtString, ok := data.(map[string]interface{})["created_at"].(string)
		if !ok {
			w.writeError("Missing 'created_at'", http.StatusBadRequest)
			return
		}
		createdAt, err := time.Parse(layout, createAtString)
		if err != nil {
			w.writeError(err.Error(), http.StatusBadRequest)
			return
		}

		viewerIDString, ok := data.(map[string]interface{})["viewer_id"].(string)
		if !ok {
			w.writeError("Missing 'viewer_id'", http.StatusBadRequest)
			return
		}

		viewerID, err := uuid.Parse(viewerIDString)
		if err != nil {
			w.writeError(err.Error(), http.StatusBadRequest)
			return
		}

		content, ok := data.(map[string]interface{})["data"]
		if !ok {
			w.writeError("Missing 'data'", http.StatusBadRequest)
			return
		}
		d := models.Data{
			CreatedAt: createdAt,
			ViewerID:  viewerID,
			Data:      content.(map[string]interface{}),
		}
		dataList[i] = d
	}

	if err := c.DBController.AddData(dataList); err != nil {
		if err.Error() == dbcontrollers.ErrFailedToAddData.Error() {
			w.writeError(err.Error(), http.StatusAccepted)
			return
		}
		w.writeError(err.Error(), http.StatusInternalServerError)
		return
	}

	w.writeData("OK", http.StatusCreated)
}

func (c *RESTController) getDataByViewer(w ResponseWriter, r *Request) {
	log.Println("Getting data by viewer")
	if err := checkRequestType(GET, w, r); err != nil {
		return
	}

	viewerIDs, ok := r.URL.Query()["viewer_id"]
	if !ok || len(viewerIDs[0]) < 1 {
		w.writeError("Url Param 'email' is missing", http.StatusBadRequest)
		return
	}

	viewerID, err := uuid.Parse(viewerIDs[0])
	if err != nil {
		w.writeError(err.Error(), http.StatusBadRequest)
		return
	}

	times, ok := r.URL.Query()["interval"]
	if !ok || len(times[0]) < 1 {
		w.writeError("Url Param 'email' is missing", http.StatusBadRequest)
		return
	}
	layout := "Mon Jan 02 2006 15:04:05.0000 GMT-0700"
	interval, err := time.Parse(layout, times[0])
	if err != nil {
		w.writeError(err.Error(), http.StatusBadRequest)
		return
	}

	projectData, err := c.DBController.GetDataByViewerAndTime(&viewerID, interval)
	if err != nil {
		if err.Error() == dbcontrollers.ErrFailedToAddData.Error() {
			w.writeError(err.Error(), http.StatusAccepted)
			return
		}
		w.writeError(err.Error(), http.StatusInternalServerError)
		return
	}

	w.writeData(projectData, http.StatusCreated)

}

package server

import (
	"errors"
	"net/http"
	"time"

	"gopkg.in/mgo.v2"

	"github.com/celrenheit/trending-machine/hubspider"
	"github.com/gorilla/context"
)

var timeLayout = "02-01-2006"

func HandleSnapshots(w http.ResponseWriter, r *http.Request) {
	database, ok := context.GetOk(r, "DB")
	if !ok {
		RespondWithError(w, r, errors.New("Could'nt obtain database"), http.StatusInternalServerError)
		return
	}
	db := database.(*mgo.Database)
	params := r.URL.Query()

	d := time.Now()
	if date, ok := params["date"]; ok {
		if t, err := time.Parse(timeLayout, date[0]); err == nil {
			d = t
		}
	}

	snapshot, err := hubspider.FindSnapshotByTime(db, d)
	if err != nil {
		if err == mgo.ErrNotFound {
			RespondWithError(w, r, err, http.StatusNotFound)
			return
		}
		RespondWithError(w, r, err, http.StatusInternalServerError)
		return
	}
	snapshot.ServeJSON(w, r)
}

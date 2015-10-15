package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/context"
	"gopkg.in/mgo.v2"
)

type Response map[string]interface{}

func DBMiddleware(session *mgo.Session, DBNAME string) negroni.Handler {
	return negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		s := session.Clone()
		defer s.Close()
		context.Set(r, "dbSession", s)
		context.Set(r, "DBNAME", DBNAME)
		context.Set(r, "DB", s.DB(DBNAME))
		next(w, r)
	})
}

func RespondWithError(w http.ResponseWriter, r *http.Request, err error, code int) {
	ServeJSON(w, r, &Response{"error": err.Error()}, code)
}

func ServeJSON(w http.ResponseWriter, r *http.Request, response *Response, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		fmt.Fprint(w, "")
	}
}

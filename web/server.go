package server

import (
	"fmt"
	"log"
	"net/http"

	"gopkg.in/mgo.v2"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
)

func Launch(dbSession *mgo.Session, DBNAME, PORT string) {

	router := mux.NewRouter()
	router.HandleFunc("/api/snapshots", HandleSnapshots).Methods("GET")

	n := negroni.Classic()
	n.Use(DBMiddleware(dbSession, DBNAME))
	n.UseHandler(router)

	static := router.PathPrefix("/").Subrouter()
	static.Methods("GET").Handler(http.FileServer(http.Dir("web/public")))

	fmt.Println("Launching server at http://localhost" + PORT)
	log.Fatal(http.ListenAndServe(PORT, n))
}

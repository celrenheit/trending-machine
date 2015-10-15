package main

import (
	"encoding/json"
	"log"
	"os"
	"path"
	"time"

	"github.com/celrenheit/spider/schedulers"
	"github.com/celrenheit/trending-machine/hubspider"
	"github.com/celrenheit/trending-machine/web"
	"gopkg.in/mgo.v2"
)

var (
	PORT        = ":3000"
	MONGODB_URL = ":27017"
	DBNAME      = "ghtrending"
)

type Settings struct {
	Languages []string `json:"languages"`
}

func init() {
	if p := os.Getenv("PORT"); p != "" {
		PORT = ":" + p
	}
	if m := os.Getenv("MONGODB_URL"); m != "" {
		MONGODB_URL = m
		_, name := path.Split(MONGODB_URL)
		DBNAME = name
	}
}

func main() {

	dbSession := DBConnect(MONGODB_URL)
	err := DBEnsureIndices(dbSession, DBNAME)
	if err != nil {
		log.Println("Error settings up indices")
		log.Fatal(err)
	}

	settings, err := readConfig()
	if err != nil {
		log.Fatal(err)
	}

	go launchScheduler(dbSession.DB(DBNAME), settings)

	server.Launch(dbSession, DBNAME, PORT)
}

func launchScheduler(db *mgo.Database, settings *Settings) {

	scheduler := schedulers.NewBasicScheduler()

	scheduler.Handle(hubspider.New(db, settings.Languages)).Every(12 * time.Hour)

	log.Fatal(scheduler.Start())
}

func readConfig() (*Settings, error) {
	file, err := os.Open("settings.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var settings Settings
	err = json.NewDecoder(file).Decode(&settings)
	if err != nil {
		return nil, err
	}
	return &settings, nil
}

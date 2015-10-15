package hubspider

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var timeLayout = "02-01-2006"

type Snapshot struct {
	Date      string           `bson:"date" json:"date"`
	Languages map[string]Repos `bson:"languages" json:"languages"`
}

func NewSnapshot(langs map[string]Repos) *Snapshot {
	y, m, d := time.Now().Date()
	return &Snapshot{
		Date:      fmt.Sprintf("%d-%d-%d", d, m, y),
		Languages: langs,
	}
}

func (s *Snapshot) ToJSON() string {
	b, err := json.Marshal(s)
	if err != nil {
		return ""
	}
	return string(b)
}

func (s *Snapshot) ServeJSON(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, s.ToJSON())
}

func (s *Snapshot) Save(db *mgo.Database) error {
	return s.Upsert(db, bson.M{"$set": s})
}

func (s *Snapshot) Upsert(db *mgo.Database, query interface{}) error {
	sC := db.C("snapshots")
	_, err := sC.Upsert(bson.M{
		"date": s.Date,
	}, query)
	return err
}

func FindSnapshotByTime(db *mgo.Database, date time.Time) (*Snapshot, error) {
	sC := db.C("snapshots")
	var snap Snapshot
	err := sC.Find(bson.M{
		"date": date.Format(timeLayout),
	}).One(&snap)
	if err != nil {
		return nil, err
	}
	return &snap, nil
}

func EnsureSnapshotsIndices(s *mgo.Session, DBNAME string) error {
	sC := s.DB(DBNAME).C("snapshots")

	return sC.EnsureIndexKey("date")
}

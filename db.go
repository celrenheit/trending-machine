package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/celrenheit/trending-machine/hubspider"
	"gopkg.in/mgo.v2"
)

func DBConnect(address string) *mgo.Session {
	session, err := mgo.Dial(address)
	if err != nil {
		panic(err)
	}
	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			log.Printf("%v captured - Closing database connection\n", sig)
			session.Close()
			os.Exit(1)
		}
	}()

	return session
}

func DBEnsureIndices(s *mgo.Session, DBNAME string) error {
	if err := hubspider.EnsureSnapshotsIndices(s, DBNAME); err != nil {
		return err
	}
	return nil
}

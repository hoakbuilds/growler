package db

import (
	"log"

	"gopkg.in/mgo.v2"
)

// MongoClient is the struct that represents the server instance.
// This server is used to estabilish connection with mongodb.
type MongoClient struct {
	DBSession *mgo.Session
	DB        *mgo.Database
	Cfg       *MongoConnCfg
}

// Connect is the function used to start the server responsible for the
// mongodb connection
func (c *MongoClient) Connect() error {

	var url string

	if c.Cfg.HostAddress != "" && c.Cfg.HostPort != "" {
		url = c.Cfg.HostPort + ":" + c.Cfg.HostAddress
	} else {

		if c.Cfg.HostAddress != "" {
			url = url + c.Cfg.HostAddress + ":"
		} else {
			url = url + DefaultMongoHost + ":"
		}
		if c.Cfg.HostPort != "" {
			url = url + c.Cfg.HostPort + ":"
		} else {
			url = url + DefaultMongoPort + ":"
		}
	}

	log.Printf("[GRWLR-DBC] Trying to connect to mongo at %s", url)

	/*
	 * Connect to the server and get a database handle
	 */
	dbSession, err := mgo.Dial(url)
	if err != nil {
		return err
	}
	defer dbSession.Close()

	log.Printf("[GRWLR-DBC] Fetching mongodb Session")
	db := dbSession.DB("growler")
	c.DBSession = dbSession
	c.DB = db

	info, err := c.DB.Session.BuildInfo()
	if err != nil {
		return err
	}
	log.Printf("[GRWLR-DBC] mongodb connection successful! mongodb build: %v", info)

	return nil
}

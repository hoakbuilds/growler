package main

import (
	"sync"

	"github.com/murlokito/growler/growler/api"
	"github.com/murlokito/growler/growler/db"
)

// App represents the application structure, it encapsulates
// a Tab Controller, a HTTP/HTTPS Server with a REST API
// and the MongoDB Client.
type App struct {
	RestServer  *api.WebService
	MongoClient *db.MongoClient
	Controller  *TabController

	WaitGroup *sync.WaitGroup

	Cfg *Config
}

// StartApp is the function used to finally start the App,
// allowing channel communication between the provided shell,
// the MongoClient and the Controller.
func (a *App) StartApp() {

}

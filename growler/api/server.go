package api

import (
	"log"
	"sync"

	"github.com/gorilla/mux"
	"github.com/murlokito/growler/growler/db"
)

// WebService is the structure that defines the
// overall webservice. This includes a REST API and
// the MongoClient used for the database.
type WebService struct {
	Router    *mux.Router
	MgoClient *db.MongoClient
	WaitGroup *sync.WaitGroup

	ServChan chan interface{}

	Cfg *WebSvcCfg

	AppChan *chan int
}

// StartServer is used to start the Web Service
func (ws *WebService) StartServer() {

	if ws.Cfg.ServeTLS == true {
		/*
		 * Start the server in a goroutine
		 */
		go func() {
			ws.BuildAndServeTLS()
		}()
	} else {
		/*
		 * Start the server in a goroutine
		 */
		go func() {
			ws.BuildAndServe()
		}()
	}

	<-ws.ServChan
	log.Printf("[GRWLR-API] Terminating API Server")
}

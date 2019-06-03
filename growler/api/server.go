package api

import (
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
}

// StartServer is used to start the Web Service
func (w *WebService) StartServer() error {

	if w.Cfg.ServeTLS == true {
		err := generateCertificate()
		if err != nil {
			return err
		}
		/*
		 * Start the server in a goroutine
		 */
		go func() {
			w.BuildAndServeTLS()
		}()
	} else {
		/*
		 * Start the server in a goroutine
		 */
		go func() {
			w.BuildAndServe()
		}()
	}

	<-w.ServChan
	return nil
}

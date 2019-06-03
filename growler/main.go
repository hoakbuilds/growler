package main

import (
	"io"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/murlokito/growler/growler/api"

	"github.com/murlokito/growler/growler/db"
)

const (
	noFileOrDirectory string = " no such file or directory"
)

var (
	wg sync.WaitGroup
)

// processInput will receive a command to execute and will perform logic prior to
// executing those commands, so that they are properly executed and terminated
// when it's needed to
func processInput(input string) []string {

	log.Print("[processInput]\n")
	// Remove the newline character.
	input = strings.TrimSuffix(input, "\n")

	inputArr := strings.Split(input, " ")

	return inputArr
}

func growlerMain(cfg Config) error {
	var (
		tabCntl      TabController
		mongoConnCfg *db.MongoConnCfg
	)
	defer wg.Wait()

	log.Printf("[GRWLR] Creating mongodb connection config")

	if cfg.mongoAddr != "" && cfg.mongoPort != "" {
		mongoConnCfg = &db.MongoConnCfg{
			HostAddress: cfg.mongoAddr,
			HostPort:    cfg.mongoPort,
		}
	} else {
		if cfg.mongoAddr != "" {
			mongoConnCfg = &db.MongoConnCfg{
				HostAddress: cfg.mongoAddr,
			}
		}
		if cfg.mongoPort != "" {
			mongoConnCfg = &db.MongoConnCfg{
				HostPort: cfg.mongoPort,
			}
		}
	}
	// try to connect to mongodb
	mongoClient := db.MongoClient{
		Cfg: mongoConnCfg,
	}

	err := mongoClient.Connect()

	if err != nil {
		return err
	}

	log.Printf("[GRWLR] Connected to mongodb at %s", cfg.mongoPort)

	log.Printf("[GRWLR] Creating REST API server config")
	apiCfg := api.WebSvcCfg{
		ServeTLS: cfg.tls,
		RestPort: cfg.restPort,
	}

	api := api.WebService{
		MgoClient: &mongoClient,
		WaitGroup: &wg,
		Cfg:       &apiCfg,
		ServChan:  make(chan interface{}),
	}

	err = api.StartServer()

	if err != nil {
		return err
	}
	// before performing a normal start, check for a previous session
	restore, latestID, err := restoreSession()

	if err != nil {
		//assume there is no previous session
		errSplit := strings.Split(err.Error(), ":")
		if errSplit[1] == noFileOrDirectory {
			log.Printf("[restoreSession]: no previous session found")
		} else {
			log.Printf("[restoreSession]: %v\n", err)
		}
	}

	tabCntlChan := make(chan string, 10)
	shellCommChan := make(chan string, 10)

	if len(restore) != 0 {
		tabCntl = TabController{
			tcChan: tabCntlChan,
			shChan: shellCommChan,
			tabs:   []Tab{},
			id:     latestID,
		}

	} else {
		tabCntl = TabController{
			tcChan: tabCntlChan,
			shChan: shellCommChan,
			tabs:   []Tab{},
			id:     0,
		}
	}

	param := tabCntl.Run
	wg.Add(1)
	go param(restore)

	f, err := os.OpenFile("info.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	mw := io.MultiWriter(os.Stdout, f)
	log.SetOutput(mw)
	log.Println("[GRWLR] Logger started.")

	app := App{
		Controller:  &tabCntl,
		RestServer:  &api,
		MongoClient: &mongoClient,
		WaitGroup:   &wg,
		Cfg:         &cfg,
	}

	app.Run()

	return nil
}

// Main is the main function of the program
func Main() {

	cfg := loadConfig()

	// Work around defer not working after os.Exit()
	if err := growlerMain(cfg); err != nil {
		log.Printf("[GRWLR] err: %v", err)
		return
	}

	return
}

func main() {
	Main()
}

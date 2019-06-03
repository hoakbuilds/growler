package main

import (
	"io"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

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

	log.Printf("[GRWLR] Connected to mongodb at %s:%s", cfg.mongoAddr, cfg.mongoPort)

	log.Printf("[GRWLR] Creating REST API server config")
	apiCfg := api.WebSvcCfg{
		ServeTLS: cfg.tls,
		RestPort: cfg.restPort,
	}

	AppChan := make(chan int, 5)
	api := api.WebService{
		MgoClient: &mongoClient,
		WaitGroup: &wg,
		Cfg:       &apiCfg,
		ServChan:  make(chan interface{}),
		AppChan:   &AppChan,
	}
	log.Printf("[GRWLR] Starting API Server..")

	apiServParam := api.StartServer
	wg.Add(1)
	go apiServParam()

	log.Printf("[GRWLR] Setting up Tab Controller")

	if err != nil {
		return err
	}
	// before performing a normal start, check for a previous session
	restore, latestID, err := restoreSession()

	if err != nil {
		//assume there is no previous session
		errSplit := strings.Split(err.Error(), ":")
		if errSplit[1] == noFileOrDirectory {
			log.Printf("[GRWLR]: no previous session found")
		} else {
			log.Printf("[GRWLR]: %v\n", err)
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

	tabCntlParam := tabCntl.Run
	wg.Add(1)
	go tabCntlParam(restore)

	f, err := os.OpenFile("info.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	mw := io.MultiWriter(os.Stdout, f)
	log.SetOutput(mw)
	log.Println("[GRWLR] Logger started.")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Printf("[GRWLR] Catching signal, terminating gracefully.")
		if len(tabCntl.tabs) > 1 {
			err := tabCntl.TerminateGracefully()
			if err != nil {
				log.Printf("[GRWLR] Error terminating tab: %v", err)
			}
		}
		os.Exit(1)
	}()
	wg.Add(1)
	app := App{
		Controller:  &tabCntl,
		RestServer:  &api,
		MongoClient: &mongoClient,
		WaitGroup:   &wg,
		Cfg:         &cfg,
		AppChan:     make(chan int, 5),
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

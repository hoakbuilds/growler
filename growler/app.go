package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

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

	AppChan chan int
}

// StartApp is the function used to finally start the App,
// allowing channel communication between the provided shell,
// the MongoClient and the Controller.
func (a *App) StartApp() {

}

// Run is a shell abstraction to get commands from the cli.
func (a *App) Run() {
	time.Sleep(500000)

	reader := bufio.NewReader(os.Stdin)

	for {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		wg.Add(1)
		go func() {
			select {
			case <-a.AppChan:
				fmt.Print("\n[GRWLR] > ")
			case <-c:
				log.Printf("\n[GRWLR] Catching signal, terminating gracefully.")
			}
		}()
		log.Print("[GRWLR] growler shell\n")
		fmt.Print("\n[GRWLR] > ")
		// Read the keyboad input.
		in, err := reader.ReadString('\n')

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		input := processInput(in)

		if a.Cfg.debug == true {
			log.Printf("[shell_debug]: %v\n", input)
		}

		if input[0] == "help" {
			fmt.Fprintf(os.Stdout, "\n \t Command usage help\n")
			fmt.Fprintf(os.Stdout, "\nget <url> \t\t\t\t\t - Opens a new tab and gives it the job of requesting the given url.\n")
			fmt.Fprintf(os.Stdout, "\nlist \t\t\t\t\t\t - Lists all open tabs, their current or last job, and their current status.\n")
			fmt.Fprintf(os.Stdout, "\nread <image_filename> \t\t\t\t - Reads an image and tries to decode it into a txt with RGB values.\n")
			fmt.Fprintf(os.Stdout, "\ndownload <url_to_download> <output_filename> \t - Requests the URL and saves the response or downloadable file to the output file.\n")
			fmt.Fprintf(os.Stdout, "\nstop <tab id> \t\t\t\t\t - Stops the thread represented by the input id.\n\t\t\t\t\t\t\t If a thread is working it will wait until it has finished to stop it.\n")
			fmt.Fprintf(os.Stdout, "\njob <url> \t\t\t\t\t - Adds the given url to the pool of jobs, it will be executed once a tab is available.\n")
			fmt.Fprintf(os.Stdout, "\nexit \t\t\t\t\t\t - Stops all tabs and exits the browser.\n\n")
			continue
		} else if input[0] == "list" {

			a.Controller.tcChan <- input[0]
			recv := <-a.Controller.shChan

			if recv == "ok" {
				continue
			}

		} else if input[0] == "get" {

			if len(input) < 2 {
				fmt.Print("[error] you must provide a valid URL.\n")
				continue
			} else {
				a.Controller.tcChan <- input[0]
				recv := <-a.Controller.shChan

				if recv == "url" {
					a.Controller.tcChan <- input[1]
				}

				ok := <-a.Controller.shChan

				if ok == "ok" {
					continue
				}

			}

		} else if input[0] == "get" && input[1] == "image" {

		} else if input[0] == "exit" {
			a.Controller.tcChan <- input[0]
			msg := <-a.Controller.shChan

			if msg == "ok" {
				break
			}

		}
		if input[0] == "stop" {

			if len(input) < 2 {
				fmt.Print("[error] you must provide a valid Tab ID.\n")
				continue
			} else {
				a.Controller.tcChan <- input[0]
				recv := <-a.Controller.shChan

				if recv == "TabID" {
					a.Controller.tcChan <- input[1]
				}

				continue
			}

		}
		if input[0] == "stop" && input[1] == "all" {

			a.Controller.tcChan <- input[0]
			recv := <-a.Controller.shChan

			if recv == "TabID" {
				a.Controller.tcChan <- input[1]
			}
			ok := <-a.Controller.shChan
			if ok == "ok" {
				continue
			}

		}

		if input[0] == "job" {

			if err != nil {
				log.Printf("[err] %s", err)
			}

		}

		if input[0] == "read" {
			if len(input) < 2 {
				fmt.Print("[error] you must provide a valid image file name.\n")
			} else {

				pixels, err := ReadImage(input[1])
				if err != nil {
					// replace this with real error handling
					panic(err.Error())
				}

				fmt.Println("Print pixels? Y/n")
				// Read the keyboad input.
				in, err := reader.ReadString('\n')

				if err != nil {
					fmt.Fprintln(os.Stderr, err)
				}

				secPrompt := processInput(in)

				filename := strings.Split(input[1], ".")
				f, err := os.OpenFile(filename[0]+".txt", os.O_RDWR|os.O_CREATE, 0666)

				if secPrompt[0] == "Y" || secPrompt[0] == "y" {

					for row := 0; row < len(pixels); row++ {

						for x := 0; x < len(pixels[row]); x++ {
							str := fmt.Sprintf("X: %d Y: %d\t(R,G,B,A) (%d, %d, %d, %d)\n", x, row, pixels[row][x].R, pixels[row][x].G, pixels[row][x].B, pixels[row][x].A)
							f.WriteString(str)
						}
					}
				}
				f.Close()
			}

		}
		if input[0] == "crawl" {
			if len(input) < 2 {
				fmt.Print("[error] you must provide a search term to crawl unsplash.com\n")
			} else {

			}
		}

		if input[0] == "download" {
			if len(input) < 3 {
				fmt.Print("[error] you must provide a valid URL to download a file and the file name to save it.\n")
			} else {
				err := DownloadFile(input[2], input[1])
				if err != nil {
					// replace this with real error handling
					panic(err.Error())
				}
			}
		}

	}

}

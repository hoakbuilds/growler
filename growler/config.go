package main

import (
	"flag"
	"os"

	"github.com/murlokito/growler/growler/db"

	"github.com/murlokito/growler/growler/api"
)

// Config is used as a structure to hold information
// passed via the CLI upon GCD startup. It's passed
// into the Server structure to define some of it's
// parameters
type Config struct {
	tls       bool
	debug     bool
	restPort  string
	mongoPort string
	mongoAddr string
}

func loadConfig() Config {

	var (
		tlsVar       string
		debugVar     string
		restPortVar  string
		mongoPortVar string
		mongoIPVar   string
	)

	flag.StringVar(&tlsVar, "tls", "", "If it is desired to use HTTPS to serve the API. Set `true` to enable.")
	flag.StringVar(&debugVar, "debug", "", "If it is desired to enable debug mode. Set `true` to enable.")
	flag.StringVar(&restPortVar, "rest", api.DefaultRestPort, "The port desired to run the REST API.")
	flag.StringVar(&mongoPortVar, "mgoaddr", "localhost", "The host address of the mongodb server.")
	flag.StringVar(&mongoIPVar, "mgoport", db.DefaultMongoPort, "The port desired to run the REST API.")

	flag.Parse()
	if len(os.Args) == 0 {
		return Config{}
	}

	if tlsVar == "true" {
		if debugVar == "true" {
			return Config{
				tls:       true,
				debug:     true,
				restPort:  restPortVar,
				mongoAddr: mongoIPVar,
				mongoPort: mongoPortVar,
			}
		}
		return Config{
			tls:      true,
			debug:    false,
			restPort: restPortVar,

			mongoAddr: mongoIPVar,
			mongoPort: mongoPortVar,
		}
	}
	return Config{
		tls:      false,
		debug:    false,
		restPort: restPortVar,

		mongoAddr: mongoIPVar,
		mongoPort: mongoPortVar,
	}

}

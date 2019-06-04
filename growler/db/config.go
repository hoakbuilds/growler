package db

const (
	// DefaultMongoHost is the default port used for the mongodb server
	// when no input port is given
	DefaultMongoHost = "localhost"
	// DefaultMongoPort is the default port used for the mongodb server
	// when no input port is given
	DefaultMongoPort = "27017"
)

// MongoConnCfg is used as a structure to hold config information
// passed via the CLI upon growler startup. It's passed
// into the WebService structure to define some of it's
// parameters
type MongoConnCfg struct {
	HostAddress string
	HostPort    string
}

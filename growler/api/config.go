package api

const (
	// DefaultRestPort is the default port used for the REST API
	// when no input port is given
	DefaultRestPort = "9000"
	// DefaultMongoURI is used as a default unchanging part that is
	// concatenated with the host address and port to make the mongodb
	// connection URL.
	DefaultMongoURI = "mongodb://"
)

// WebSvcCfg is used as a structure to hold config information
// passed via the CLI upon growler startup. It's passed
// into the WebService structure to define some of it's
// parameters
type WebSvcCfg struct {
	ServeTLS bool
	RestPort string
}

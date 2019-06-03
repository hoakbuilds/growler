package api

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// Route is a structure that defines the endpoints of the API.
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// Routes is a type that encapsulates a list of Route elements
type Routes []Route

// NewRouter creates a new router based on Gorilla's mux router
// It also wraps the handlers with logging functionality.
func (ws *WebService) NewRouter() *mux.Router {

	var routes = Routes{
		Route{
			"Index",
			"GET",
			"/",
			ws.Index,
		},
		Route{
			"ImageStream",
			"GET",
			"/images/stream/{id}",
			ws.GetImages,
		},
	}

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler

		handler = route.HandlerFunc
		handler = ws.Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}

// Logger is a wrapper for the http handler.
// It gets passed the handler and returns the same handler
// with added logging and timing functionalities.
func (ws *WebService) Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		log.Printf(
			"[GRWLR-API] %s\t%s\t%s\t%s",
			r.Method,
			r.RequestURI,
			name,
			time.Since(start),
		)
		*ws.AppChan <- 1
	})
}

// BuildAndServe is the function used to serve the API endpoints
func (ws *WebService) BuildAndServe() {
	log.Println("[GRLWR-API] Building API endpoints.")
	*ws.AppChan <- 1
	ws.Router = ws.NewRouter()

	log.Printf("[GRLWR-API] Listening and Serving API. Port: %v", ws.Cfg.RestPort)
	*ws.AppChan <- 1
	log.Fatal(http.ListenAndServe(":"+ws.Cfg.RestPort,
		handlers.CORS(
			handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
			handlers.AllowedMethods([]string{"GET", "POST"}),
			handlers.AllowedOrigins([]string{"*"}),
		)(ws.Router),
	),
	)

}

// BuildAndServeTLS is the function used to serve the API endpoints
func (ws *WebService) BuildAndServeTLS() {
	*ws.AppChan <- 1
	log.Println("[GRLWR-API] Building API endpoints.")

	ws.Router = ws.NewRouter()
	*ws.AppChan <- 1
	log.Printf("[GRLWR-API] Listening and Serving API. Port: %v", ws.Cfg.RestPort)

	log.Fatal(http.ListenAndServeTLS(":"+ws.Cfg.RestPort, "cert.pem", "private.key",
		handlers.CORS(
			handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
			handlers.AllowedMethods([]string{"GET", "POST"}),
			handlers.AllowedOrigins([]string{"*"}),
		)(ws.Router),
	),
	)
}

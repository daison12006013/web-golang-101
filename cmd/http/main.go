package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"regexp"

	"web-golang-101/pkg/utils"

	_ "web-golang-101/docs"

	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

var port string
var dsn string
var appKey string

func init() {
	flag.StringVar(&port, "port", "8080", "Port to run the server on")
	flag.StringVar(&dsn, "dsn", "", "Sentry DSN")
	flag.StringVar(&appKey, "appKey", "", "Application Key")
	flag.Parse()

	port = utils.GetEnvWithDefault("APP_PORT", port)
	utils.InitializeSentry(utils.GetEnvWithDefault("SENTRY_DSN", dsn))
	utils.InitializeAppKey(utils.GetEnvWithDefault("APP_KEY", appKey))
}

// @title Web Golang 101 API
func main() {
	r := chi.NewRouter()

	apiRouter := ApiRoutes()

	localhostRouter := chi.NewRouter()
	localhostRouter.Mount("/api", apiRouter)

	hr := &HostRouter{}
	hr.Map("^localhost:\\d+$", localhostRouter)
	hr.Map("api\\.(.*)", apiRouter)
	r.Mount("/", hr)

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	fmt.Printf("Server is listening on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

type Route struct {
	Pattern *regexp.Regexp
	Handler http.Handler
}

type HostRouter struct {
	Routes []*Route
}

func (hr *HostRouter) Map(pattern string, handler http.Handler) {
	re := regexp.MustCompile(pattern)
	route := &Route{Pattern: re, Handler: handler}
	hr.Routes = append(hr.Routes, route)
}

func (hr *HostRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, route := range hr.Routes {
		if route.Pattern.MatchString(r.Host) {
			route.Handler.ServeHTTP(w, r)
			return
		}
	}
	http.NotFound(w, r)
}

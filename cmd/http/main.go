package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"regexp"

	"web-golang-101/pkg/utils"

	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
)

var port string
var dsn string
var appKey string

func init() {
	flag.StringVar(&port, "port", "8080", "Port to run the server on")
	flag.StringVar(&dsn, "dsn", "", "Sentry DSN")
	flag.StringVar(&appKey, "appKey", "", "Application Key")
	flag.Parse()

	port = utils.GetEnvWithDefault("PORT", port)
	utils.InitializeSentry(utils.GetEnvWithDefault("SENTRY_DSN", dsn))
	utils.InitializeAppKey(utils.GetEnvWithDefault("APP_KEY", appKey))
}

func main() {
	r := chi.NewRouter()

	utils.SetDefaultMiddlewares(r)

	apiRouter := ApiRoutes()
	webRouter := WebRoutes()
	adminRouter := AdminRoutes()

	localhostRouter := chi.NewRouter()
	localhostRouter.Mount("/api", apiRouter)
	localhostRouter.Mount("/admin", adminRouter)
	localhostRouter.Mount("/", webRouter)

	hr := &HostRouter{}
	hr.Map("^localhost:\\d+$", localhostRouter)
	hr.Map("api\\.(.*)", apiRouter)
	hr.Map("admin\\.(.*)", adminRouter)
	hr.Map("(.*)", webRouter)
	r.Mount("/", hr)

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

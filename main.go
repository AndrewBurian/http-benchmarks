package main

import (
	"net/http"
	"github.com/andrewburian/http-benchmarks/routes"
)


type Router struct {
	Name        string
	SetupRoutes func(routes []routes.Route, param string) http.Handler
}

func main() {

	// create the list of routers
	routers := []Router{
		{"Powermux", SetupRoutesPowermux},
	}

	RunTests(routers)
}

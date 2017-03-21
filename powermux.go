package main

import (
	"net/http"
	"io"
	"github.com/andrewburian/powermux"
	"github.com/andrewburian/http-benchmarks/routes"
)

type powermuxHandler string

func (p powermuxHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, powermux.GetPathParam(r, string(p)))
}

func SetupRoutesPowermux(routes []routes.Route, param string) http.Handler {
	h := powermuxHandler(param)
	r := powermux.NewServeMux()
	for _, route := range routes {
		r.Route(route.Path).Any(h)
	}
	return r
}

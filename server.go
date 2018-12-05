package main

import (
	"fmt"
	"net/http"

	"github.com/ServiceComputingTeam/Blog-Server/jsonp"
	"github.com/ServiceComputingTeam/Blog-Server/jwt"
	"github.com/urfave/negroni"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "Welcome to the home page!")
	})
	mux.HandleFunc("/login", func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		username := req.Form.Get("username")
		password := req.Form.Get("password")
		if username == "user" && password == "123456" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
		}
	})

	n := negroni.New()
	n.Use(jsonp.NewJsonp())
	n.Use(jwt.NewJwt())
	// n.Use(negroni.HandlerFunc(jwt.ValidatorJWT))
	n.UseHandler(mux)

	n.Run(":3000")
}

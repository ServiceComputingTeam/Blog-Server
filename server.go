package main

import (
	sw "github.com/ServiceComputingTeam/Blog-Server/go"
	"github.com/ServiceComputingTeam/Blog-Server/jsonp"
	"github.com/ServiceComputingTeam/Blog-Server/jwt"
	"github.com/urfave/negroni"
)

func main() {

	router := sw.NewRouter()
	userRouter := sw.NewUserRouter()

	router.PathPrefix("/user").Handler(negroni.New(
		jwt.NewJwt(),
		negroni.HandlerFunc(jwt.ValidatorJWT),
		negroni.Wrap(userRouter),
	))

	n := negroni.New()
	n.Use(jsonp.NewJsonp())
	n.UseHandler(router)
	n.Run(":3000")
}

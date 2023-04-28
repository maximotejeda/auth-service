package main

import (
	"os"

	r "github.com/maximotejeda/auth-service/external/v0.0.1/router"
)

var (
	port = os.Getenv("SERVERPORT")
	addr = os.Getenv("SERVERADDR")
)

func main() {

	router := r.NewRouter()

	r.AuthAddRoutes(router)
	r.UserAddRoutes(router)
	r.AdminAddRoutes(router)

	r.R.Run(addr + ":" + port)
}

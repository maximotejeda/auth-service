package main

import (
	r "github.com/maximotejeda/auth-service/external/v0.1/router"
)

func main() {
	router := r.NewRouter()

	r.AuthAddRoutes(router)
	r.UserAddRoutes(router)
	r.AdminAddRoutes(router)

	r.R.Run(":8083")
}

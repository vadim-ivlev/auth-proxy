package main

import (
	"auth-proxy/router"
)

func main() {
	router.Serve(":4000")
}

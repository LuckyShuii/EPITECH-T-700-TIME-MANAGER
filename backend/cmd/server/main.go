package main

import (
	"app/internal/router"
)

func main() {
	r := router.SetupRouter()
	r.Run(":5000")
}
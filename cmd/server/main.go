package main

import (
	"vm-hub/internal/infrustructure/server"
)

func main() {
	r := server.SetupRouter()
	server := server.NewServer(&r)
	server.Start()
}

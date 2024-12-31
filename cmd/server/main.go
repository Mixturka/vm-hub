package main

import (
	"vm-hub/internal/config"
	"vm-hub/internal/infrustructure/server"
)

func main() {
	config := config.LoadConfig()

	r := server.SetupRouter()
	server := server.NewServer(&r)
	server.Start(config)
}

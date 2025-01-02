package main

import (
	"log/slog"
	"os"

	"github.com/Mixturka/vm-hub/internal/config"
	"github.com/Mixturka/vm-hub/internal/infrustructure/server"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	r := server.SetupRouter()
	server := server.NewServer(&r)
	server.Start(config)
}

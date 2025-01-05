package main

import (
	"log/slog"
	"os"

	"github.com/Mixturka/vm-hub/internal/app/config"
	"github.com/Mixturka/vm-hub/internal/app/infrustructure/server"
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

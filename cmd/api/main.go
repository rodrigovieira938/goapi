package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/rodrigovieira938/goapi/api/router"
	"github.com/rodrigovieira938/goapi/config"
	"github.com/rodrigovieira938/goapi/util"
	"github.com/rodrigovieira938/goapi/util/logger"
)

func main() {
	logger.Init()
	cfg := config.New()
	if _, err := os.Stat(".env"); err != nil {
		fmt.Println("Environment file not found!\nCreating .env")
		err = config.WriteToFile(cfg, ".env")
		if err != nil {
			fmt.Println("Error creating .env file:", err)
		} else {
			fmt.Println(".env file created successfully.")
		}
	} else {
		config.DebugPrint(cfg)
	}
	slog.Info("Starting server...")
	http.Handle("/", router.New())
	port := cfg.Server.Port
	if port == 0 {
		port = util.FindUsablePort(8080)
	}
	slog.Info(fmt.Sprintf("Server started at http://%s:%d", cfg.Server.Hostname, port))
	err := http.ListenAndServe(fmt.Sprintf("%s:%d", cfg.Server.Hostname, port), nil)
	if err != nil {
		slog.Error("Error starting server:", "err", err)
	}
}

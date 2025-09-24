package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/rodrigovieira938/goapi/api/router"
	"github.com/rodrigovieira938/goapi/config"
	"github.com/rodrigovieira938/goapi/util"
)

func main() {
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

	http.Handle("/", router.New())
	port := cfg.Server.Port
	if port == 0 {
		port = util.FindUsablePort(8080)
	}
	fmt.Printf("Server started at http://%s:%d\n", cfg.Server.Hostname, port)
	err := http.ListenAndServe(fmt.Sprintf("%s:%d", cfg.Server.Hostname, port), nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}

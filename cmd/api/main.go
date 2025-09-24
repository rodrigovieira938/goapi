package main

import (
	"fmt"
	"net/http"

	"github.com/rodrigovieira938/goapi/api/router"
	"github.com/rodrigovieira938/goapi/util"
)

func main() {
	http.Handle("/", router.New())
	port := util.FindUsablePort(8080)
	fmt.Printf("Server started at http://localhost:%d\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}

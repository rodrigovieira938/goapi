package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/rodrigovieira938/goapi/config"
	"github.com/rodrigovieira938/goapi/util/db"
)

func usage() {
	fmt.Printf("Usage: %s [command] [flags]\n", flag.CommandLine.Name())
	fmt.Printf("\tup - Run the migration\n")
	fmt.Printf("\tdown - Rollback the migration\n")
}

func migrate_up() {
	fmt.Println("Migrating up...")
	cfg := config.New()
	conn, err := db.Connect(cfg.Database)
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}
	defer conn.Close()

	entries, err := os.ReadDir("./migrations/up")
	if err != nil {
		log.Fatal(err)
	}

	for _, e := range entries {
		content, err := os.ReadFile("./migrations/up/" + e.Name())
		if err != nil {
			log.Fatal("Failed to read migration file "+e.Name()+":", err)
		}
		str := string(content)
		_, err = conn.Exec(str)
		if err != nil {
			log.Fatal("Failed to execute migration "+e.Name()+":", err)
		}
	}
}
func migrate_down() {
	fmt.Println("Migrating down...")
	cfg := config.New()
	conn, err := db.Connect(cfg.Database)
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}
	defer conn.Close()

	entries, err := os.ReadDir("./migrations/down")
	if err != nil {
		log.Fatal(err)
	}

	for _, e := range entries {
		content, err := os.ReadFile("./migrations/down/" + e.Name())
		if err != nil {
			log.Fatal("Failed to read migration file "+e.Name()+":", err)
		}
		str := string(content)
		_, err = conn.Exec(str)
		if err != nil {
			log.Fatal("Failed to execute migration "+e.Name()+":", err)
		}
	}
}

func main() {
	flag.Parse() // Make sure there is no flags since there is none
	if len(flag.Args()) == 0 || flag.Arg(0) != "up" && flag.Arg(0) != "down" {
		usage()
		return
	}
	if flag.Arg(0) == "up" {
		migrate_up()
	} else if flag.Arg(0) == "down" {
		migrate_down()
	}
}

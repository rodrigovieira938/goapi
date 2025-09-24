package main

import (
	"flag"
	"fmt"
)

func usage() {
	fmt.Printf("Usage: %s [command] [flags]\n", flag.CommandLine.Name())
	fmt.Printf("\tup - Run the migration\n")
	fmt.Printf("\tdown - Rollback the migration\n")
}

func migrate_up() {
	//TODO: implement
	fmt.Println("Migrating up...")
}
func migrate_down() {
	//TODO: implement
	fmt.Println("Migrating down...")
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

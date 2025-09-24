package config

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/joeshaw/envdecode"
	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
}
type ServerConfig struct {
	Hostname string `env:"SERVER_HOSTNAME,default=localhost"`
	Port     int    `env:"SERVER_PORT,default=8080"`
	Debug    bool   `env:"SERVER_DEBUG,default=false"`
}
type DatabaseConfig struct {
	Hostname string `env:"DB_HOST,default=localhost"`
	Port     int    `env:"DB_PORT,default=5432"`
	User     string `env:"DB_USER,default=postgres"`
	Password string `env:"DB_PASSWORD,default=password"`
}

func New() *Config {
	var c Config
	if err := envdecode.StrictDecode(&c); err != nil {
		slog.Error("Failed to decode", "err", err)
	}
	return &c
}
func WriteToFile(c *Config, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	config, err := envdecode.Export(c)
	if err != nil {
		return err
	}
	for _, v := range config {
		fmt.Fprintf(file, "%s=%s\n", v.EnvVar, v.Value)
	}
	return nil
}
func DebugPrint(c *Config) {
	fmt.Println("Current Environment:")
	fmt.Println("\tServer Configuration:")
	fmt.Printf("\t\tHostname=%s\n", c.Server.Hostname)
	fmt.Printf("\t\tPort=%d\n", c.Server.Port)
	fmt.Printf("\t\tDebug=%t\n", c.Server.Debug)
	fmt.Println("\tDatabase Configuration:")
	fmt.Printf("\t\tHostname=%s\n", c.Database.Hostname)
	fmt.Printf("\t\tPort=%d\n", c.Database.Port)
	fmt.Printf("\t\tUser=%s\n", c.Database.User)
	fmt.Printf("\t\tPassword=%s\n", strings.Repeat("*", len(c.Database.Password)))
}

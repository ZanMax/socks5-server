package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/armon/go-socks5"
)

type Config struct {
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Address  string `json:"address"`
}

func loadConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %v", err)
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("failed to decode config: %v", err)
	}
	return &config, nil
}

func main() {
	// Load configuration
	config, err := loadConfig("config.json")
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	// Create credentials
	credentials := socks5.StaticCredentials{
		config.Username: config.Password,
	}
	authenticator := socks5.UserPassAuthenticator{Credentials: credentials}

	// Create SOCKS5 server configuration
	conf := &socks5.Config{
		AuthMethods: []socks5.Authenticator{authenticator},
	}

	// Create SOCKS5 server
	server, err := socks5.New(conf)
	if err != nil {
		fmt.Printf("Error creating server: %v\n", err)
		os.Exit(1)
	}

	// Start server
	listenAddr := fmt.Sprintf("%s:%d", config.Address, config.Port)
	fmt.Printf("Starting SOCKS5 server on %s\n", listenAddr)
	if err := server.ListenAndServe("tcp", listenAddr); err != nil {
		fmt.Printf("Error running server: %v\n", err)
		os.Exit(1)
	}
}

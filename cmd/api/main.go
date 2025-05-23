package main

import (
	"context"
	"lingo/utils"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// Load config
	config, err := utils.LoadConfig("../../")
	if err != nil {
		log.Fatalf("Couldn't fetch config due to %s", err.Error())
	}

	// Connect to database
	pgConfig, err := pgxpool.ParseConfig(config.DBSource)
	if err != nil {
		log.Fatalf("Couldn't parse config due to %s", err.Error())
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), pgConfig)
	if err != nil {
		log.Fatalf("Couldn't connect to database due to %s", err.Error())
	}
	defer pool.Close()

	// Initialize server
	server := NewServer(pool, config)
	if err := server.router.Run(config.ServerAddr); err != nil {
		log.Fatalf("Couldn't start server due to %s", err.Error())
	}
	log.Println("Server started on port", config.ServerAddr)
}

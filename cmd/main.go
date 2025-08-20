package main

import (
	"blog-api/config"
	"blog-api/internal/server"
	"log"
)

func main() {
	if err := config.Init(); err != nil {
		log.Fatalf("init config err: %v", err)
	}
	cfg := config.GetConfig()

	s, err := server.NewServer(cfg)
	if err != nil {
		log.Fatalf("failed to init server: %v", err)
	}

	if err := s.Run(); err != nil {
		log.Fatalf("server stopped with error: %v", err)
	}
}

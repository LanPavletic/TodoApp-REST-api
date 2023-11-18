package main

import (
	"log"

	"github.com/LanPavletic/go-rest-server/internal/server"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	var err error
	s := server.New()
	s.RegisterRoutes()
	s.Start()
	err = s.Stop()
	if err != nil {
		log.Fatal("There was a problem closing the server: ", err)
	}
}

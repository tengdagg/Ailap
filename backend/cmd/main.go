package main

import (
	"log"
	"os"

	"ailap-backend/internal/server"
)

func main() {
	if err := server.StartHTTPServer(); err != nil {
		log.Println("server exit with error:", err)
		os.Exit(1)
	}
}





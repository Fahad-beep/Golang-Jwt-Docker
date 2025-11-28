package main

import (
	"log"
	"main/internal/config"
	"main/internal/routes"
	"net/http"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env")
	}

	db := config.ConnectionDB()
	r := routes.SetupRouter(db)

	log.Println("Running on port 8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal("error running on port 8080", err)
	}
}

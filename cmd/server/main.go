package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ayushvyasonwork/chatapp1/internal/db"
	"github.com/ayushvyasonwork/chatapp1/internal/middlewares"
	"github.com/ayushvyasonwork/chatapp1/internal/routes"
	"github.com/joho/godotenv"
)

func main() {
	// load .env if present (ignore error if file missing so env can be provided differently)
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found, falling back to environment variables")
	}

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("MONGO_URI not set")
	}

	if err := db.ConnectMongo(mongoURI); err != nil {
		log.Fatal("Mongo connection failed:", err)
	}

	routes.RegisterRoutes()

	fmt.Println("🚀 Go Chat Server started on :8000")
	handler := middlewares.CORSMiddleware(http.DefaultServeMux)

	log.Fatal(http.ListenAndServe(":8000", handler))
}

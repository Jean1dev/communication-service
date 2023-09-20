package main

import (
	"communication-service/infra/sockets"
	"communication-service/routes"
	"log"
	"net/http"
	"os"

	"github.com/getsentry/sentry-go"
	sentryhttp "github.com/getsentry/sentry-go/http"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		os.Setenv("MONGO_URI", "mongodb://localhost:27017")
		log.Printf("Error loading .env file %s", err.Error())
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Print("server running on ", port)

	setupAPI()
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Panic(err)
	}

}

func setupAPI() {
	if err := sentry.Init(sentry.ClientOptions{
		Dsn:              "https://533bc6327d19b7a619643db76175d214@o318666.ingest.sentry.io/4505874771804160",
		EnableTracing:    true,
		TracesSampleRate: 1.0,
	}); err != nil {
		log.Fatalf("Sentry initialization failed: %v", err)
	}

	sentryHandler := sentryhttp.New(sentryhttp.Options{})
	http.HandleFunc("/email", sentryHandler.HandleFunc(routes.EmailHandler))
	socketsManager := sockets.NewManager()
	http.HandleFunc("/ws", socketsManager.ServeWS)
}

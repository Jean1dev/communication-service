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
	dir, _ := os.Getwd()
	log.Printf("Diretorio atual %v", dir)

	err := godotenv.Load()
	if err != nil {
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
	http.HandleFunc("/email-stats", sentryHandler.HandleFunc(routes.EmailEstatisticasHandler))

	http.HandleFunc("/notificacao", sentryHandler.HandleFunc(routes.NotificationHandler))
	http.HandleFunc("/notificacao/mark-as-read", sentryHandler.HandleFunc(routes.NotificationHandler))

	http.HandleFunc("/social-feed", sentryHandler.HandleFunc(routes.SocialFeedHandler))
	http.HandleFunc("/social-feed/", sentryHandler.HandleFunc(routes.SocialFeedHandler))
	http.HandleFunc("/social-feed/like", sentryHandler.HandleFunc(routes.SocialFeedHandler))
	http.HandleFunc("/social-feed/comment", sentryHandler.HandleFunc(routes.SocialFeedHandler))
	socketsManager := sockets.InitManagerGlobally()
	http.HandleFunc("/ws", socketsManager.ServeWS)
}

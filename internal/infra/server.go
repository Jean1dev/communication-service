package infra

import (
	"log"
	"net/http"
	"os"

	"github.com/Jean1dev/communication-service/api"
	"github.com/Jean1dev/communication-service/internal/infra/sockets"
	"github.com/getsentry/sentry-go"
	sentryhttp "github.com/getsentry/sentry-go/http"
)

func ConfigAndStartHttpServer() {
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
	http.HandleFunc("/email", sentryHandler.HandleFunc(api.EmailHandler))
	http.HandleFunc("/email-stats", sentryHandler.HandleFunc(api.EmailEstatisticasHandler))

	http.HandleFunc("/notificacao", sentryHandler.HandleFunc(api.NotificationHandler))
	http.HandleFunc("/notificacao/mark-as-read", sentryHandler.HandleFunc(api.NotificationHandler))

	http.HandleFunc("/social-feed", sentryHandler.HandleFunc(api.SocialFeedHandler))
	http.HandleFunc("/social-feed/", sentryHandler.HandleFunc(api.SocialFeedHandler))
	http.HandleFunc("/social-feed/like", sentryHandler.HandleFunc(api.SocialFeedHandler))
	http.HandleFunc("/social-feed/comment", sentryHandler.HandleFunc(api.SocialFeedHandler))
	socketsManager := sockets.InitManagerGlobally()
	http.HandleFunc("/ws", socketsManager.ServeWS)
}

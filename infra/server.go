package main

import (
	"communication-service/routes"
	"fmt"
	"net/http"
	"os"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Print(w, "Hi")
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Print("server running on ", port)
	http.HandleFunc("/email", routes.EmailHandler)
	http.ListenAndServe(":"+port, nil)
}

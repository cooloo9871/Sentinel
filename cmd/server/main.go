package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/brobridge/sentinel/internal/handler"
	k8sclient "github.com/brobridge/sentinel/internal/k8s"
)

func main() {
	dynClient, err := k8sclient.NewDynamicClient()
	if err != nil {
		log.Fatalf("k8s client: %v", err)
	}

	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	if len(jwtSecret) == 0 {
		log.Fatal("JWT_SECRET environment variable is required")
	}

	namespace := k8sclient.CurrentNamespace()
	store := k8sclient.NewStore(dynClient)

	cfg := handler.Config{
		Store:     store,
		DynClient: dynClient,
		JWTSecret: jwtSecret,
		Namespace: namespace,
		TokenTTL:  8 * time.Hour,
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("starting sentinel on :%s (namespace: %s)", port, namespace)
	if err := http.ListenAndServe(":"+port, handler.New(cfg)); err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"io/fs"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/brobridge/sentinel/internal/handler"
	k8sclient "github.com/brobridge/sentinel/internal/k8s"
	sentinelweb "github.com/brobridge/sentinel/web"
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

	mux := http.NewServeMux()
	mux.Handle("/api/", handler.New(cfg))

	staticFS, err := fs.Sub(sentinelweb.StaticFiles, "dist")
	if err != nil {
		log.Fatalf("embed sub: %v", err)
	}
	mux.Handle("/", spaHandler(http.FS(staticFS)))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("starting sentinel on :%s (namespace: %s)", port, namespace)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}

// spaHandler serves static files and falls back to index.html for SPA routing.
func spaHandler(fsys http.FileSystem) http.Handler {
	fileServer := http.FileServer(fsys)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f, err := fsys.Open(r.URL.Path)
		if err != nil {
			r.URL.Path = "/"
		} else {
			f.Close()
		}
		fileServer.ServeHTTP(w, r)
	})
}

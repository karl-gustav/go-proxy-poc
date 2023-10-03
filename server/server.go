package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()
	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		j, err := json.Marshal(map[string]any{
			"url": r.URL,
		})
		if err != nil {
			panic(err)
		}
		w.Write(j)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	log.Println("Serving http://localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

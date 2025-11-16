package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func CreateAndRunServer(r *chi.Mux, port string, errChan chan error) *http.Server {
	serv := http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: r,
	}

	go func() {
		defer log.Printf("server stop working...\n")

		log.Printf("server run on %s\n", port)
		if err := serv.ListenAndServe(); err != nil {
			log.Printf("server error: %v\n", err)
			errChan <- err
			return
		}
	}()

	return &serv
}

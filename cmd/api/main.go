package main 

import ( 
	"fmt"
	"net/http"
	"github.com/sam8beard/csv-json-api/internal/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() { 
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Post("/upload", handlers.UploadHandler)

	server := &http.Server{
		Addr: ":8080",
		Handler: router,
	} // server

	fmt.Println("Listening on port", server.Addr)
	err := server.ListenAndServe(); if err != nil {fmt.Println("Failed to listen to server", err)}

} // main
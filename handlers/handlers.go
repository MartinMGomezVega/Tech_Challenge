package handlers

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

// drivers: Setear el puerto y escuchar el servidor (se encarga de manejar las solicitudes HTTP que llegan al servidor)
func Drivers() {
	router := mux.NewRouter() // Devuelve informacion del router

	// abrir el puerto
	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "8080"
	}

	// Creando un handler/importacion
	// Los cors son quienes otorgan los permisos
	handler := cors.AllowAll().Handler(router)        // Todos pueden acceder
	log.Fatal(http.ListenAndServe(":"+PORT, handler)) // Escucha el puerto
}

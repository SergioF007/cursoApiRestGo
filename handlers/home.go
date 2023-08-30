package handlers

import (
	"net/http"

	"platzi.com/go/api/rest-ws/server"
)

// Creamos la el struct de la respuesta que se le devolvera al cliente
type HomeResponse struct {
	Message string `json:"message"`
	Status  bool   `json:"status"`
}

func HomeHandler(s server.Server) http.HandlerFunc {
	// w es la respuesta que le hacemos llegar al cliente y r es la data que nos envia el cliente
	return func(w http.ResponseWriter, r *http.Request) {
	
}
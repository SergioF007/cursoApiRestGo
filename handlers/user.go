package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/segmentio/ksuid"
	"platzi.com/go/api/rest-ws/models"
	"platzi.com/go/api/rest-ws/repository"
	"platzi.com/go/api/rest-ws/server"
)

// esto son los datos necesarios para que un usuario sea capaz de registrarse dentro de la aplicacion de
type SigUpRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Definimos el struct de de la respuesta cuando se registra
type SigUpResponse struct {
	Id    string `json:"id"`
	Email string `json:"email"`
}

// Controladora de registro
func SignUpHandler(s server.Server) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		var request = SigUpRequest{}
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		id, err := ksuid.NewRandom()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var user = models.User{
			Email:    request.Email,
			Password: request.Password,
			Id:       id.String(),
		}
		err = repository.InsertUser(r.Context(), &user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(SigUpResponse{
			Id:    user.Id,
			Email: user.Email,
		})

	}
}

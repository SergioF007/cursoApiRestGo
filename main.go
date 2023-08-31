package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"platzi.com/go/api/rest-ws/handlers"
	"platzi.com/go/api/rest-ws/server"
)

func main() {

	// cargamos el paquete de las variables de entorno
	err := godotenv.Load(".env")
	if err != nil {
		panic("Error loading.env file")
	}

	// las capturamos en esta variables
	PORT := os.Getenv("PORT")
	JWT_SECRET := os.Getenv("JWT_SECRET")
	DATABASE_URL := os.Getenv("DATABASE_URL")

	// creamos la conexion asignado los datos a a las variables correcpondiente del servidor
	s, err := server.NewServer(context.Background(), &server.Config{
		Port:        PORT,
		JWTSecret:   JWT_SECRET,
		DatabaseUrl: DATABASE_URL,
	})

	if err != nil {
		log.Fatal(err)
	}

	s.Start(BindRoutes)

}

// para poder poder poner a corrrer la conexion con la funcion Start, necesitamos el binder
// por eso segun la estrucra que se espera creamos esta funcion
func BindRoutes(s server.Server, r *mux.Router) {
	r.HandleFunc("/", handlers.HomeHandler(s)).Methods(http.MethodGet)
	r.HandleFunc("/signup", handlers.SignUpHandler(s)).Methods(http.MethodPost)
}

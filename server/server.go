package server

import (
	"context"
	"errors"

	"github.com/gorilla/mux"
)

type Config struct {
	Port        string // el puerto donde se va a ejecutar
	JWTSecret   string // una llave secreta para generar los tokens
	DatabaseUrl string // conexion a la base de datos
}

type Server interface {
	Config() *Config
}

type Broker struct {
	config *Config
	router *mux.Router
}

// crea la funcion para hacer que el struc se comporte como tipo server
func (b *Broker) Config() *Config {
	return b.config
}

// Constructor para el struct Broker que recibe la funcion config
func NewServer(ctx context.Context, config *Config) (*Broker, error) {

	// revisar la config para asegurar de que no tenga campos vacios
	if config.Port == "" {
		return nil, errors.New("port is required") // mesg: el puerto es requerido
	}

	// replicamos este comportamiento para los demas

	if config.JWTSecret == "" {
		return nil, errors.New("secret is required") // mesg: el puerto es requerido
	}

	if config.DatabaseUrl == "" {
		return nil, errors.New("database url is required")
	}

	return &Broker{
		config: config,
		router: mux.NewRouter(),
	}, nil
}

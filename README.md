# cursoApiRestGo

# Curso Go Api-Rest

## ****CRUD****

- **C**reate
- **R**ead
- **U**pdate
- **D**elete

A nivel de HTTP se crea una URL, esto sirve como referencia para un recurso por ejemplo`/posts/`.

Usando los diferentes métodos de HTTP indicamos la operación que deseamos realizar:

Create -> **POST** `/posts`

Read -> **GET** `/posts`

Read -> **GET** `/posts/:id`

Update -> **PUT** `/posts/:id`

Delete -> **DELETE** `/posts/:id`

Creamos el archivo go.mod del proyecto con: 

Darle permisos de usuario al paquete donde estoy ubicado: 

`sudo chown -R ingsergio:ingsergio /home/ingsergio/go/src/rest-websockets/`

luego si ejecutamos el siguiente comando:

```go
go mod init [platzi.com/go/api/rest-ws](http://platzi.com/go/api/rest-ws) 
```

**Intenté hacerlo agregando a la ruta la v1, pero no lo pudo crear por un mal proceder, que se explica a continuación**. 

El error que estás viendo se debe a que estás intentando inicializar un módulo Go con una versión en la ruta del módulo (**`v1`** en este caso), pero Go solo permite esto para versiones v2 o posteriores.

Si estás trabajando en la primera versión de tu módulo, simplemente omite la versión en la ruta del módulo:

```bash
bashCopy code
go mod init platzi.com/go/api/rest-ws

```

Si estás trabajando en la versión 2 o posterior, entonces el sufijo de versión es obligatorio y debe ser **`/v2`**, **`/v3`**, etc.:

```bash
bashCopy code
go mod init platzi.com/go/api/rest-ws/v2

```

Recuerda que el sufijo de versión en la ruta del módulo es una característica específica de Go Modules y es obligatorio para importar versiones v2 o posteriores de un módulo. Para la versión 1, simplemente no se utiliza un sufijo de versión.

Descargamos 4 dependencias principales:

```go
go get github.com/gorilla/mux  // modulo para tener un router y un websockets
go get github.com/golang-jwt/jwt/v4   // modulo que nos ayuda a crear web tokens 
go get github.com/gorilla/websocket  // modulo para el manejo de los websockets
go get github.com/joho/godotenv   // para manejar la variables de entorno
```

podemos ver nuestra dependencias instaladas en el archivo .mod

![Untitled](https://s3-us-west-2.amazonaws.com/secure.notion-static.com/31c06dd8-5dee-4326-af62-62ef267c94a3/Untitled.png)

Crear Carpeta server y el paquete server.go

## ****Struct Server****

Creamos el struc del Broker para manejar los servidores, se crear la función para poder heredar el comportamiento de config que tiene el struc del servido y se crea el constructor que garantiza la configuración que se debe de tener cumpla con los datos necesarios para crear la conexión a la base de datos

```go
package server

import (
	"context"
	"errors"
	"log"
	"net/http"

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

// Creamos la funcion que se encarga le levantarce o ejecutarce
func (b *Broker) Start(binder func(s Server, r *mux.Router)) {
	b.router = mux.NewRouter() // para crear nuevos router
	binder(b, b.router)
	log.Println("Starting server on port", b.Config().Port)
	if err := http.ListenAndServe(b.Config().Port, b.router); err != nil {
		log.Fatal("ListenAndServe:", err)
	}

}
```

## ****Nuestro primer endpoint****

************************************************************Configuración de variables de entorno (Archvio .env)************************************************************

1. Es el por donde va a correr todo 
2. La llave de nuestro Tokens - Secret
3. La conexión a la Base de datos 

![Untitled](https://s3-us-west-2.amazonaws.com/secure.notion-static.com/feb2fb5b-ac17-4a4c-ae7b-c21c7ca66ed0/Untitled.png)

**Se creo el archivo .env para la variables de entorno y se construyo la cabecera de la respuesta que va a recibir el cliente cuando su petición de conexión fue exitosa**

![Untitled](https://s3-us-west-2.amazonaws.com/secure.notion-static.com/ffb10849-d1d4-4066-9ea4-451987c46b56/Untitled.png)

```go
package handlers

import (
	"encoding/json"
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
		w.Header().Set("Content-Type", "application/json") // le definimos el header a la peticion 
		// aqui le idicamos al cliente que nostros le estamos enviando es un JSON 
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(HomeResponse{
			Message: "Welcome to Platzi GO",
			Status:  true, 
		})
	}
}
```

**********************************************En el Archivo main cargamos el archivó de las variables de entorno  y inicializamos las variables de entorno para hacer la conexión con el servidor de la BD, Creamos la función  BindRoutes que es el tipo de método que espera la función Start que es la que maneja la conexión.**********************************************  

```go
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
}
```

## ****Definiendo el modelo de usuarios y Persistencia****

```go
package models

type User struct {
	Id       int64  `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
```

## ****Patrón repository****

El erro al hacer una implementación **concretas** a una base de datos, que luego si necesitó estableces conexión con otra base de datos vamos a empezar a volver el código demasiado volátil y mas complicado de manejar.   Por eso es mejor implementar abstracciones. 

![Untitled](https://s3-us-west-2.amazonaws.com/secure.notion-static.com/c92edcf6-c9f3-4599-a320-028420f4f59b/Untitled.png)

![Untitled](https://s3-us-west-2.amazonaws.com/secure.notion-static.com/40c12130-eb4a-492e-83e5-bbdfc769fe96/Untitled.png)

```go
// Concretas 
Handler - GetUserByIdPostgreSQL 

// Absatracciones
Handler - GetUserById - User
			Postgres
			MongoDB
			...
```

```go
package repository

import (
	"context"

	"platzi.com/go/api/rest-ws/models"
)

type UserRepository interface {
	InsertUser(ctx context.Context, user *models.User) error
	GetUserById(ctx context.Context, id int64) (*models.User, error)
	Close() error
}

var implementation UserRepository // este UserRepository puede ser una implementacion en Mongo o MySWL o PostgreSQL

func SetRepository(repository UserRepository) {
	implementation = repository
}

func InsertUser(ctx context.Context, user *models.User) error {
	return implementation.InsertUser(ctx, user)
}

func GetUserById(ctx context.Context, id int64) (*models.User, error) {
	return implementation.GetUserById(ctx, id)
}

func Close() error {
	return implementation.Close()
}
```

**Nota: Este un patrón bastante potente ya que nos permite implementar métodos flexibles para crear abstracciones de nuestro código sin preocuparse lo la base de datos.** 

## ****Registro de usuarios****

para este registro de usuario vamos a trabajar con PostgreSql así que creamos el repo con los métodos que vamos a implementar para sus respectivas **consulta** a la **Base de datos.**  

```go
package database

import (
	"context"
	"database/sql"
	"log"

	"platzi.com/go/api/rest-ws/models"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(url string) (*PostgresRepository, error) {

	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	return &PostgresRepository{db}, nil
}

func (repo *PostgresRepository) InsertUser(ctx context.Context, user *models.User) error {
	_, err := repo.db.ExecContext(ctx, "INSERT INTO user (email, password) VALUES ($1, $2)", user.Email, user.Password)
	return err
}

func (repo *PostgresRepository) GetUserById(ctx context.Context, id int64) (*models.User, error) {
	rows, err := repo.db.QueryContext(ctx, "SELECT id, email FROM users WHERE id = $1", id)
	if err != nil {
		log.Fatal(err.Error())
		return nil, err
	}
	defer rows.Close()

	var user = models.User{}
	for rows.Next() {
		if err := rows.Scan(&user.Id, &user.Email); err != nil {
			return &user, nil
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &user, nil

	//Nota:  esta parte puedo mejorar, tomar como ejemplo la forma en lo que lo hacen Javier en vedana

}

func (repo *PostgresRepository) Close() error {
	return repo.db.Close()
}
```

## ****Implementando el registro****

Vamos a utilizar Doker como una herramienta para levantar nuestra base de datos.  

Vamos a descargar una librería que implementa muchas de las cosas que PostgreSQL utiliza en Go. 

```go
go get github.com/lib/pq  
```

![Untitled](https://s3-us-west-2.amazonaws.com/secure.notion-static.com/ef41d2f1-12af-4f0c-a990-138837e52e06/Untitled.png)

Vamos a crear el DockerFile  nos va ayudar a levantar una instancia de PostgreSQL sin tener que instalarlo. 

Agregamos una librería para que nos genere texto aleatorio 

```go
go get github.com/segmentio/ksuid 
```

![Untitled](https://s3-us-west-2.amazonaws.com/secure.notion-static.com/d1afc932-851c-4213-b2ac-0cb3ed3aeecf/Untitled.png)

![Untitled](https://s3-us-west-2.amazonaws.com/secure.notion-static.com/130ee5a9-f2f4-4b81-96d7-2e2d7a66d11b/Untitled.png)

**en handlers creamos los struc necesarios para el registro de usuario como el struc de su respuesta con el metodo que se encarga de hacer el registro**

![Untitled](https://s3-us-west-2.amazonaws.com/secure.notion-static.com/98213b28-9208-493f-a769-aca57675dbbc/Untitled.png)

```go
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
```

## ****Probando los registros****

![Untitled](https://s3-us-west-2.amazonaws.com/secure.notion-static.com/1ead64ca-f7b4-4886-b4c1-f77d3ba3a279/Untitled.png)

![Untitled](https://s3-us-west-2.amazonaws.com/secure.notion-static.com/b25c90b2-1802-4b47-9f0f-89d84007655f/Untitled.png)

Una forma mas facil de crear la base de datos por docker-compose:

```
version: "3.8"

services:
  db:
    image: postgres
    restart: always
    container_name: database
    ports:
      - "5432:5432"
    environment:
      POSTGRES_USER: root
      POSTGRES_PASSWORD: root
      POSTGRES_DB: pgdb
      PGDATA: /var/lib/postgresql/data
    volumes:
      - ./db-scripts:/docker-entrypoint-initdb.d

```

donde db-scripts es la ruta donde se guarden los scripts para inicializar la base de datos

tener en cuenta que el schema en este caso se llama pgdb

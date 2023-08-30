package server

type Config struct {
	Port        string // el puerto donde se va a ejecutar
	JWTSecret   string // una llave secreta para generar los tokens
	DatabaseUrl string // conexion a la base de datos
}

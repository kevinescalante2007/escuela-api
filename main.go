package main

import (
	"escuela_api/config"
	"escuela_api/routes"
	//"github.com/gin-gonic/gin"
)

func main() {
	// Conectar a la base de datos
	config.ConnectDB()

	// Configurar y arrancar el servidor
	r := routes.SetupRouter()
	r.Run("192.168.1.7:8080") // El servidor estará corriendo en http://localhost:8080
	//r.Run("10.79.16.28:8080") // El servidor estará corriendo en http://localhost:8080
}

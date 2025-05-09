package config

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql" // Importante para la conexión MySQL
)

var DB *sql.DB

func ConnectDB() {
	var err error
	// Configuración de conexión: usuario, contraseña, dirección (localhost:3306) y base de datos
	DB, err = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/escuela")
	//DB, err = sql.Open("mysql", "root:@tcp(10.79.14.151:3306)/escuela")
	if err != nil {
		fmt.Println("Error al conectar a la base de datos:", err)
		return
	}

	// Verificar la conexión
	err = DB.Ping()
	if err != nil {
		fmt.Println("Error al hacer ping a la base de datos:", err)
		return
	}

	fmt.Println("Conexión a la base de datos exitosa")
}

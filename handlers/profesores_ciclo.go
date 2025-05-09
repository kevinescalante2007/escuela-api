package handlers

import (
	//"database/sql"
	"database/sql"
	"escuela_api/config"
	"escuela_api/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Obtener todos los registros de Profesor_Ciclo
func GetProfesorCiclos(c *gin.Context) {
	rows, err := config.DB.Query("SELECT Ciclo, CI_Profesor FROM Profesor_Ciclo")
	if err != nil {
		// Aquí agregamos un log detallado del error para rastrear el problema
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al consultar la base de datos", "details": err.Error()})
		return
	}
	defer rows.Close()

	var ciclos []models.ProfesorCiclo
	for rows.Next() {
		var ciclo models.ProfesorCiclo
		if err := rows.Scan(&ciclo.Ciclo, &ciclo.CIProfesor); err != nil {
			// Log detallado de error al leer datos
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al leer datos", "details": err.Error()})
			return
		}
		ciclos = append(ciclos, ciclo)
	}

	c.JSON(http.StatusOK, ciclos)
}


// Crear un registro en Profesor_Ciclo
func CreateProfesorCiclo(c *gin.Context) {
	var ciclo models.ProfesorCiclo
	if err := c.ShouldBindJSON(&ciclo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos incorrectos"})
		return
	}

	// 🚨 VALIDACIÓN que necesitas agregar:
	if ciclo.Ciclo == "" || ciclo.CIProfesor == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ciclo y CI_Profesor no pueden estar vacíos"})
		return
	}

	_, err := config.DB.Exec("INSERT INTO Profesor_Ciclo (Ciclo, CI_Profesor) VALUES (?, ?)", ciclo.Ciclo, ciclo.CIProfesor)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al insertar en la base de datos"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Profesor asignado al ciclo correctamente"})
}


// Eliminar una asignación de profesor a un ciclo
func DeleteProfesorCiclo(c *gin.Context) {
	ciclo := c.Param("ciclo")
	ciProfesor := c.Param("ci")

	res, err := config.DB.Exec("DELETE FROM Profesor_Ciclo WHERE Ciclo = ? AND CI_Profesor = ?", ciclo, ciProfesor)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar la asignación"})
		return
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Relación profesor-ciclo no encontrada"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Asignación eliminada correctamente"})
}





// GetProfesorCicloByCI obtiene un ciclo de profesor por su CI
func GetProfesorCicloByCI(c *gin.Context) {
	ci := c.Param("ci")
	var profesorCiclo models.ProfesorCiclo

	// Consulta en la base de datos
	err := config.DB.QueryRow("SELECT Ciclo, CI_Profesor FROM Profesor_Ciclo WHERE CI_Profesor = ?", ci).Scan(&profesorCiclo.Ciclo, &profesorCiclo.CIProfesor)
	if err != nil {
		// Si no encuentra el ciclo o ocurre un error
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Profesor Ciclo no encontrado"})
		} else {
			// Aquí agregamos un log detallado del error para rastrear el problema
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al consultar la base de datos", "details": err.Error()})
		}
		return
	}

	// Si todo sale bien, se retorna el ciclo de profesor
	c.JSON(http.StatusOK, profesorCiclo)
}

// UpdateProfesorCiclo actualiza la información de un ciclo de profesor
func UpdateProfesorCiclo(c *gin.Context) {
	ci := c.Param("ci")
	var ciclo models.ProfesorCiclo
	if err := c.ShouldBindJSON(&ciclo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos incorrectos"})
		return
	}

	// Actualiza el ciclo del profesor en la base de datos
	_, err := config.DB.Exec("UPDATE Profesor_Ciclo SET Ciclo = ? WHERE CI_Profesor = ?", ciclo.Ciclo, ci)
	if err != nil {
		// Si ocurre un error al realizar la actualización
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar el ciclo del profesor"})
		return
	}

	// Si todo salió bien, se responde con un mensaje
	c.JSON(http.StatusOK, gin.H{"message": "Ciclo del profesor actualizado correctamente"})
}

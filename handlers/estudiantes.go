package handlers

import (
	"escuela_api/config"
	"escuela_api/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Obtener todos los estudiantes
func GetEstudiantes(c *gin.Context) {
	rows, err := config.DB.Query("SELECT CI_Estudiante, Nombres FROM estudiantes")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al consultar estudiantes", "details": err.Error()})
		return
	}
	defer rows.Close()

	var estudiantes []models.Estudiante
	for rows.Next() {
		var e models.Estudiante
		if err := rows.Scan(&e.CI, &e.Nombres); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al leer estudiante", "details": err.Error()})
			return
		}
		estudiantes = append(estudiantes, e)
	}

	c.JSON(http.StatusOK, estudiantes)
}

// Crear estudiante
func CreateEstudiante(c *gin.Context) {
	var e models.Estudiante
	if err := c.ShouldBindJSON(&e); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos", "details": err.Error()})
		return
	}

	var exists bool
	err := config.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM estudiantes WHERE CI_Estudiante = ?)", e.CI).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al verificar existencia del estudiante", "details": err.Error()})
		return
	}

	if exists {
		c.JSON(http.StatusConflict, gin.H{"message": "Estudiante ya existe"})
		return
	}

	_, err = config.DB.Exec("INSERT INTO estudiantes (CI_Estudiante, Nombres) VALUES (?, ?)", e.CI, e.Nombres)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al insertar estudiante", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Estudiante creado exitosamente"})
}

// Obtener estudiante por CI
func GetEstudianteByID(c *gin.Context) {
	ci := c.Param("ci")

	var e models.Estudiante
	err := config.DB.QueryRow("SELECT CI_Estudiante, Nombres FROM estudiantes WHERE CI_Estudiante = ?", ci).
		Scan(&e.CI, &e.Nombres)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Estudiante no encontrado", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, e)
}

// Actualizar estudiante
func UpdateEstudiante(c *gin.Context) {
	ci := c.Param("ci")
	var e models.Estudiante

	if err := c.ShouldBindJSON(&e); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos", "details": err.Error()})
		return
	}

	var count int
	err := config.DB.QueryRow("SELECT COUNT(*) FROM estudiantes WHERE CI_Estudiante = ?", ci).Scan(&count)
	if err != nil || count == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Estudiante no encontrado"})
		return
	}

	_, err = config.DB.Exec("UPDATE estudiantes SET Nombres = ? WHERE CI_Estudiante = ?", e.Nombres, ci)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar estudiante", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Estudiante actualizado"})
}

// Eliminar estudiante
func DeleteEstudiante(c *gin.Context) {
	ci := c.Param("ci")

	var count int
	err := config.DB.QueryRow("SELECT COUNT(*) FROM estudiantes WHERE CI_Estudiante = ?", ci).Scan(&count)
	if err != nil || count == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Estudiante no encontrado"})
		return
	}

	// Eliminar notas asociadas
	_, err = config.DB.Exec("DELETE FROM notas WHERE CI_Estudiante = ?", ci)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar notas asociadas", "details": err.Error()})
		return
	}

	// Eliminar estudiante
	_, err = config.DB.Exec("DELETE FROM estudiantes WHERE CI_Estudiante = ?", ci)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar estudiante", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Estudiante y sus notas eliminados"})
}

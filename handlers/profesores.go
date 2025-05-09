package handlers

import (
	"database/sql"
	"escuela_api/config"
	"escuela_api/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetProfesores(c *gin.Context) {
	rows, err := config.DB.Query("SELECT ci, nombres FROM profesores")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener profesores"})
		return
	}
	defer rows.Close()

	var profesores []models.Profesor
	for rows.Next() {
		var p models.Profesor
		if err := rows.Scan(&p.CI, &p.Nombres); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al escanear profesores"})
			return
		}
		profesores = append(profesores, p)
	}
	c.JSON(http.StatusOK, profesores)
}

func GetProfesorByCI(c *gin.Context) {
	ci := c.Param("ci")
	var p models.Profesor
	err := config.DB.QueryRow("SELECT ci, nombres FROM profesores WHERE ci = ?", ci).Scan(&p.CI, &p.Nombres)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Profesor no encontrado"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener profesor"})
		return
	}
	c.JSON(http.StatusOK, p)
}

func CreateProfesor(c *gin.Context) {
	var p models.Profesor
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	_, err := config.DB.Exec("INSERT INTO profesores (ci, nombres) VALUES (?, ?)", p.CI, p.Nombres)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear profesor"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Profesor creado"})
}

func UpdateProfesor(c *gin.Context) {
	ci := c.Param("ci")
	var p models.Profesor
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	_, err := config.DB.Exec("UPDATE profesores SET nombres = ? WHERE ci = ?", p.Nombres, ci)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar profesor"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Profesor actualizado"})
}

func DeleteProfesor(c *gin.Context) {
	ci := c.Param("ci")
	_, err := config.DB.Exec("DELETE FROM profesores WHERE ci = ?", ci)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar profesor"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Profesor eliminado"})
}


package handlers

import (
	"escuela_api/config"
	"escuela_api/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Obtener todas las asignaturas
func GetAsignaturas(c *gin.Context) {
	rows, err := config.DB.Query("SELECT Id_Asignatura, Asignatura, Nivel FROM asignaturas")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al consultar asignaturas", "details": err.Error()})
		return
	}
	defer rows.Close()

	var asignaturas []models.Asignatura
	for rows.Next() {
		var a models.Asignatura
		if err := rows.Scan(&a.Id_Asignatura, &a.Asignatura, &a.Nivel); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al leer asignatura"})
			return
		}
		asignaturas = append(asignaturas, a)
	}

	c.JSON(http.StatusOK, asignaturas)
}

// Obtener una asignatura por ID
func GetAsignaturaByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var a models.Asignatura
	err = config.DB.QueryRow("SELECT Id_Asignatura, Asignatura, Nivel FROM asignaturas WHERE Id_Asignatura = ?", id).
		Scan(&a.Id_Asignatura, &a.Asignatura, &a.Nivel)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Asignatura no encontrada"})
		return
	}

	c.JSON(http.StatusOK, a)
}

// Crear o actualizar asignatura
func CreateAsignatura(c *gin.Context) {
	var a models.Asignatura

	if err := c.ShouldBindJSON(&a); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido"})
		return
	}

	// Verificar si ya existe
	var exists bool
	err := config.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM asignaturas WHERE id_asignatura = ?)", a.Id_Asignatura).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al verificar existencia"})
		return
	}

	if exists {
		c.JSON(http.StatusConflict, gin.H{"error": "Asignatura ya existe"})
		return
	}

	// Insertar nueva asignatura
	_, err = config.DB.Exec(
		"INSERT INTO asignaturas (id_asignatura, asignatura, nivel) VALUES (?, ?, ?)",
		a.Id_Asignatura, a.Asignatura, a.Nivel,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear asignatura"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Asignatura creada correctamente"})
}


// Actualizar asignatura
func UpdateAsignatura(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var a models.Asignatura
	if err := c.ShouldBindJSON(&a); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	result, err := config.DB.Exec("UPDATE asignaturas SET Asignatura = ?, Nivel = ? WHERE Id_Asignatura = ?", a.Asignatura, a.Nivel, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar asignatura"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Asignatura no encontrada"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Asignatura actualizada correctamente"})
}

// Eliminar asignatura
func DeleteAsignatura(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	result, err := config.DB.Exec("DELETE FROM asignaturas WHERE Id_Asignatura = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar asignatura"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Asignatura no encontrada"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Asignatura eliminada correctamente"})
}


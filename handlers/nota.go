package handlers

import (
	"database/sql"
	"escuela_api/config"
	"escuela_api/models"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// func calcularSupletorio(n1, n2 float64) float64 {
// 	promedio := (n1 + n2) / 2
// 	if promedio < 7 {
// 		return 0 // valor por defecto, luego puede actualizarse
// 	}
// 	return 0
// }

// GET /notas
func GetNotas(c *gin.Context) {
	rows, err := config.DB.Query("SELECT id_nota, ci_estudiante, id_asignatura, n1, n2, supletorio, updated_at FROM notas")
	if err != nil {
		log.Println("❌ Error al obtener las notas:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener las notas"})
		return
	}
	defer rows.Close()

	notas := make([]models.Nota, 0)

	for rows.Next() {
		var nota models.Nota
		var updatedAtRaw sql.NullString
		var supletorioRaw sql.NullFloat64

		err := rows.Scan(
			&nota.IdNota,
			&nota.CIEstudiante,
			&nota.IdAsignatura,
			&nota.N1,
			&nota.N2,
			&supletorioRaw,
			&updatedAtRaw,
		)
		if err != nil {
			log.Println("❌ Error escaneando fila de nota:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al procesar una nota"})
			return
		}

		if supletorioRaw.Valid {
			nota.Supletorio = &supletorioRaw.Float64
		} else {
			nota.Supletorio = nil
		}

		if updatedAtRaw.Valid {
			parsedTime, err := time.Parse("2006-01-02 15:04:05", updatedAtRaw.String)
			if err != nil {
				log.Println("⚠️ Error parseando updated_at:", updatedAtRaw.String, err)
				nota.UpdatedAt = nil
			} else {
				nota.UpdatedAt = &parsedTime
			}
		}

		notas = append(notas, nota)
	}

	c.JSON(http.StatusOK, notas)
}

// GET /notas/:id
func GetNotaByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var n models.Nota
	var supletorio sql.NullFloat64
	var updatedAt sql.NullTime

	err = config.DB.QueryRow("SELECT id_nota, ci_estudiante, id_asignatura, n1, n2, supletorio, updated_at FROM notas WHERE id_nota = ?", id).
		Scan(&n.IdNota, &n.CIEstudiante, &n.IdAsignatura, &n.N1, &n.N2, &supletorio, &updatedAt)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Nota no encontrada"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al consultar nota", "details": err.Error()})
		return
	}

	if supletorio.Valid {
		n.Supletorio = &supletorio.Float64
	}
	if updatedAt.Valid {
		t := updatedAt.Time.UTC()
		n.UpdatedAt = &t
	}

	c.JSON(http.StatusOK, n)
}

// POST /notas
func CreateNota(c *gin.Context) {
	var nota models.Nota

	if err := c.ShouldBindJSON(&nota); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido"})
		return
	}

	if nota.CIEstudiante == "" || nota.IdAsignatura == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "CI del estudiante y ID de asignatura son obligatorios"})
		return
	}

	var exists bool
	err := config.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM notas WHERE ci_estudiante = ? AND id_asignatura = ?)",
		nota.CIEstudiante, nota.IdAsignatura).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al verificar existencia de la nota"})
		return
	}
	if exists {
		c.JSON(http.StatusConflict, gin.H{"error": "La nota ya existe"})
		return
	}

	if nota.N1 < 0 || nota.N1 > 10 || nota.N2 < 0 || nota.N2 > 10 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Notas deben estar entre 0 y 10"})
		return
	}

	// Calcular si aplica supletorio
	if (nota.N1 < 7 || nota.N2 < 7) && nota.Supletorio != nil && *nota.Supletorio >= 7 {
		nota.N1 = 7
		nota.N2 = 7
	}

	updatedAt := time.Now().UTC()

	_, err = config.DB.Exec(
		`INSERT INTO notas (ci_estudiante, id_asignatura, n1, n2, supletorio, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?)`,

		nota.CIEstudiante, nota.IdAsignatura, nota.N1, nota.N2, nota.Supletorio, updatedAt,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear la nota"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Nota creada correctamente"})
}

// PUT /notas/:id
// PUT /notas/:id
// PUT /notas/:id
func UpdateNota(c *gin.Context) {
	id := c.Param("id") // Asegúrate de que 'id' está siendo correctamente extraído del parámetro de URL

	var input struct {
		N1         float64  `json:"n1"`
		N2         float64  `json:"n2"`
		Supletorio *float64 `json:"supletorio"`
	}

	// Validación de los datos entrantes
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error en los datos de entrada"})
		return
	}

	if input.N1 < 0 || input.N1 > 10 || input.N2 < 0 || input.N2 > 10 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Notas deben estar entre 0 y 10"})
		return
	}

	// Verificamos si la nota existe
	var exists bool
	err := config.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM notas WHERE Id_Nota = ?)", id).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al verificar existencia de la nota"})
		return
	}
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Nota no encontrada"})
		return
	}

	// Ejecutamos la actualización
	var result sql.Result
	if input.Supletorio != nil {
		result, err = config.DB.Exec(`
			UPDATE notas SET N1 = ?, N2 = ?, Supletorio = ?, updated_at = CURRENT_TIMESTAMP WHERE Id_Nota = ?
		`, input.N1, input.N2, *input.Supletorio, id)
	} else {
		result, err = config.DB.Exec(`
			UPDATE notas SET N1 = ?, N2 = ?, Supletorio = NULL, updated_at = CURRENT_TIMESTAMP WHERE Id_Nota = ?
		`, input.N1, input.N2, id)
	}

	// Si ocurrió un error
	if err != nil {
		log.Printf("💥 ERROR en UPDATE nota ID %s: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar nota", "details": err.Error()})
		return
	}

	// Verificamos si se actualizó alguna fila
	affected, _ := result.RowsAffected()
	if affected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No se actualizó ninguna fila"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Nota actualizada correctamente"})
}



// PUT /nota/compuesta/:ci/:id_asignatura
func UpdateNotaCompuesta(c *gin.Context) {
	ci := c.Param("ci")
	idAsignatura, err := strconv.Atoi(c.Param("id_asignatura"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de asignatura inválido"})
		return
	}

	var n struct {
		N1         float64 `json:"n1"`
		N2         float64 `json:"n2"`
		Supletorio float64 `json:"supletorio"`
	}

	if err := c.ShouldBindJSON(&n); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido"})
		return
	}

	if n.N1 < 0 || n.N1 > 10 || n.N2 < 0 || n.N2 > 10 || n.Supletorio < 0 || n.Supletorio > 10 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Notas deben estar entre 0 y 10"})
		return
	}

	var exists bool
	err = config.DB.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM notas WHERE CI_Estudiante = ? AND Id_Asignatura = ?)",
		ci, idAsignatura,
	).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al verificar existencia"})
		return
	}
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Nota no encontrada"})
		return
	}

	_, err = config.DB.Exec(
		`UPDATE notas SET N1 = ?, N2 = ?, Supletorio = ?, updated_at = CURRENT_TIMESTAMP WHERE CI_Estudiante = ? AND Id_Asignatura = ?`,
		n.N1, n.N2, n.Supletorio, ci, idAsignatura,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar nota compuesta"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Nota compuesta actualizada correctamente"})
}


// DELETE /nota/:id
func DeleteNota(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	result, err := config.DB.Exec("DELETE FROM notas WHERE id_nota = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar la nota"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Nota no encontrada"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Nota eliminada"})
}

// DELETE /nota/compuesta/:ci/:id_asignatura
func DeleteNotaCompuesta(c *gin.Context) {
	ci := c.Param("ci")
	idAsignatura, err := strconv.Atoi(c.Param("id_asignatura"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de asignatura inválido"})
		return
	}

	result, err := config.DB.Exec("DELETE FROM notas WHERE ci_estudiante = ? AND id_asignatura = ?", ci, idAsignatura)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar la nota"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Nota no encontrada"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Nota compuesta eliminada correctamente"})
}

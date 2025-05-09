package handlers

import (
	"fmt"
	"net/http"
	"time"
	"github.com/gin-gonic/gin"
)

// CompararFechasHandler compara las fechas de Python y Go
func CompararFechasHandler(c *gin.Context) {
	// Fecha en Go en formato YYYY-MM-DD HH:MM:SS (esto lo obtienes de tu base de datos en Go)
	fechaGo := "2025-04-22 22:30:28"

	// Parsear la fecha de Go
	layoutGo := "2006-01-02 15:04:05"
	tGo, err := time.Parse(layoutGo, fechaGo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error al parsear fecha de Go: %v", err)})
		return
	}

	// Convertir la fecha de Go a UTC
	tGoUTC := tGo.UTC()

	// Leer el cuerpo de la solicitud (que contendría la fecha de Python)
	var requestData struct {
		FechaPython string `json:"fecha_python"`
	}
	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error al leer los datos"})
		return
	}

	// Convertir la fecha de Python (en formato ISO 8601) a tipo time.Time
	tPython, err := time.Parse(time.RFC3339, requestData.FechaPython)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Error al parsear fecha de Python: %v", err)})
		return
	}

	// Convertir la fecha de Python a UTC si no lo está
	tPythonUTC := tPython.UTC()

	// Comparar las fechas
	var resultado string
	if tPythonUTC.After(tGoUTC) {
		resultado = "La fecha de Python es más reciente"
	} else if tPythonUTC.Before(tGoUTC) {
		resultado = "La fecha de Go es más reciente"
	} else {
		resultado = "Las fechas son iguales"
	}

	// Responder con el resultado
	c.JSON(http.StatusOK, gin.H{"resultado": resultado})
}

package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"escuela_api/handlers"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Middleware CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: true,
	}))

	// Estudiantes
	r.GET("/estudiantes", handlers.GetEstudiantes)
	r.POST("/estudiantes", handlers.CreateEstudiante)
	r.GET("/estudiantes/:ci", handlers.GetEstudianteByID)
	r.PUT("/estudiantes/:ci", handlers.UpdateEstudiante)
	r.DELETE("/estudiantes/:ci", handlers.DeleteEstudiante)

	// Profesores
	r.GET("/profesores", handlers.GetProfesores)
	r.POST("/profesores", handlers.CreateProfesor)
	r.GET("/profesores/:ci", handlers.GetProfesorByCI)
	r.PUT("/profesores/:ci", handlers.UpdateProfesor)
	r.DELETE("/profesores/:ci", handlers.DeleteProfesor)

	// Asignaturas
	r.GET("/asignaturas", handlers.GetAsignaturas)
	r.POST("/asignaturas", handlers.CreateAsignatura)
	r.GET("/asignaturas/:id", handlers.GetAsignaturaByID)
	r.PUT("/asignaturas/:id", handlers.UpdateAsignatura)
	r.DELETE("/asignaturas/:id", handlers.DeleteAsignatura)

	// Nota
	r.GET("/nota", handlers.GetNotas)
	r.POST("/nota", handlers.CreateNota)
	r.GET("/nota/:id", handlers.GetNotaByID)
	r.PUT("/nota/:id", handlers.UpdateNota)
	r.PUT("/nota/compuesta/:ci/:id_asignatura", handlers.UpdateNotaCompuesta)
	r.DELETE("/nota/:id", handlers.DeleteNota)
	r.DELETE("/nota/compuesta/:ci/:id_asignatura", handlers.DeleteNotaCompuesta)

	// Profesor Ciclo
	r.GET("/profesorciclo", handlers.GetProfesorCiclos)
	r.POST("/profesorciclo", handlers.CreateProfesorCiclo)
	r.DELETE("/profesorciclo/:ci/:ciclo", handlers.DeleteProfesorCiclo)
	r.GET("/profesorciclo/:ci", handlers.GetProfesorCicloByCI) // Faltaba esta ruta GET
	r.PUT("/profesorciclo/:ci", handlers.UpdateProfesorCiclo)  

	// Comparar fechas
	r.POST("/comparar-fechas", handlers.CompararFechasHandler)

	return r
}

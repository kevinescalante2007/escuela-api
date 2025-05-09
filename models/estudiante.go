package models

// Estudiante representa un estudiante en la base de datos
type Estudiante struct {
    CI      string `json:"ci"`       // VARCHAR(10), clave primaria
    Nombres string `json:"nombres"`  // VARCHAR(100)
}

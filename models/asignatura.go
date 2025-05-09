package models

// Asignatura representa una asignatura en el sistema
type Asignatura struct {
	Id_Asignatura int    `json:"id_asignatura"`
	Asignatura    string `json:"asignatura"`
	Nivel         string `json:"nivel"`
}

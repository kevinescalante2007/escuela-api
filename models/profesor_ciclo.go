package models

// ProfesorCiclo representa la relación entre un profesor y un ciclo académico
type ProfesorCiclo struct {
	Ciclo      string `json:"Ciclo"`
	CIProfesor string `json:"Ci_Profesor"`
}

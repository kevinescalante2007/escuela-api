package models

import "time"

type Nota struct {
	IdNota       int        `json:"id_nota"`
	CIEstudiante string     `json:"ci_estudiante"`
	IdAsignatura int        `json:"id_asignatura"`
	N1           float64    `json:"n1"`
	N2           float64    `json:"n2"`
	Supletorio   *float64   `json:"supletorio,omitempty"`
	UpdatedAt    *time.Time `json:"updated_at,omitempty"`
}

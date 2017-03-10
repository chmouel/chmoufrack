package server

import "database/sql"

var (
	ACL = true
	DB  *sql.DB
)

type errorUnauthorized struct {
	s string
}

func (e *errorUnauthorized) Error() string {
	return e.s
}

type error404 struct {
	s string
}

func (e *error404) Error() string {
	return e.s
}

type Step struct {
	ID         int     `json:"id"`
	Effort     string  `json:"effort,omitempty" yaml:"effort,omitempty"`
	EffortType string  `json:"effort_type,omitempty" yaml:"effort_type,omitempty"`
	Laps       int     `json:"laps,omitempty" yaml:"laps,omitempty"`
	Length     int     `json:"length,omitempty" yaml:"length,omitempty"`
	Percentage int     `json:"percentage,omitempty" yaml:"percentage,omitempty"`
	Type       string  `json:"type,omitempty" yaml:"type,omitempty"`
	Repeat     Repeats `json:"repeat,omitempty" yaml:"repeat,omitempty"`
	Rest       string  `json:"rest,omitempty" yaml:"rest,omitempty"`
	Position   int     `json:"-"`
}

type FBinfo struct {
	ID   int    `json:"id" facebook:"-"`
	Name string `json:"name"`
	Link string `json:"link"`
}

type Exercise struct {
	ID      int    `json:"id"`
	Name    string `json:"name" binding:"required"`
	Comment string `json:"comment,omitempty" yaml:"comment,omitempty"`
	Steps   `json:"steps,omitempty"`
	FB      FBinfo `json:"fb"`
}

type Repeats struct {
	ID      int   `json:"id"`
	Steps   Steps `json:"steps,omitempty" yaml:"steps,omitempty"`
	Repeats int   `json:"repeat,omitempty" yaml:"repeat,omitempty" `
}

type Steps []Step

func (slice Steps) Len() int {
	return len(slice)
}

func (slice Steps) Less(i, j int) bool {
	return slice[i].Position < slice[j].Position
}

func (slice Steps) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

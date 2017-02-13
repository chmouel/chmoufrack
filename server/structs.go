package server

type Repeat struct {
	ID       int    `json:"-"`
	Steps    []Step `json:"steps,omitempty"`
	Repeat   int    `json:"repeat,omitempty"`
	Position int    `json:"position,omitempty"`
}

type Step struct {
	Effort     string `json:"effort,omitempty"`
	EffortType string `json:"effort_type,omitempty"`
	Laps       int    `json:"laps,omitempty"`
	Length     int    `json:"length,omitempty"`
	Percentage int    `json:"percentage,omitempty"`
	Type       string `json:"type,omitempty"`
	Repeat     Repeat `json:"repeat,omitempty"`
	Rest       string `json:"rest,omitempty"`
	Position   int    `json:"position,omitempty"`
}

type Excercise struct {
	ID      int    `json:"-"`
	Name    string `json:",omitempty"`
	Comment string `json:",omitempty"`
	Steps   []Step `json:",omitempty"`
}

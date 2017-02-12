package server

type Repeat struct {
	ID     int    `json:",omitempty"`
	Steps  []Step `json:",omitempty"`
	Repeat int    `json:",omitempty"`
}

type Step struct {
	Effort     string `json:",omitempty"`
	EffortType string `json:",omitempty"`
	Laps       int    `json:",omitempty"`
	Length     int    `json:",omitempty"`
	Percentage int    `json:",omitempty"`
	Type       string `json:",omitempty"`
	Repeat     Repeat `json:",omitempty"`
	Rest       string `json:",omitempty"`
	Position   int    `json:",omitempty"`
}

type Excercise struct {
	ID      int    `json:",omitempty"`
	Name    string `json:",omitempty"`
	Comment string `json:",omitempty"`
	Steps   []Step `json:",omitempty"`
}

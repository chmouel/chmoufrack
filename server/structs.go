package server

type error404 struct {
	s string
}

func (e *error404) Error() string {
	return e.s
}

type Repeat struct {
	ID     int    `json:"id"`
	Steps  []Step `json:"steps,omitempty"`
	Repeat int    `json:"repeat,omitempty"`
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
	Position   int    `json:"-"`
}

type Excercise struct {
	ID      int    `json:"id"`
	Name    string `json:"name,omitempty"`
	Comment string `json:"comment,omitempty"`
	Steps   `json:"steps,omitempty"`
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

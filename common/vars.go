package common

var (
	Yaml_File = `frack.yaml`
	StaticDir = `./ui/`
)

type Workout struct {
	Repetition string `yaml:"laps"`
	Meters     string `yaml:"length"`
	Percentage string `yaml:"percentage"`
	Repos      string `yaml:"rest"`
}

type Program struct {
	Name     string
	Comment  string
	Workouts []Workout
}

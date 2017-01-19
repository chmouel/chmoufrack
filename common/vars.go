package common

import "time"

var (
	Yaml_File   = `frack.yaml`
	TrackLength = 400
	TargetVma   = "14:19"
)

type Workout struct {
	Repetition  string `yaml:"laps"`
	Meters      string `yaml:"length"`
	Percentage  string `yaml:"percentage"`
	TrackLaps   string //Temporary for convenience
	TrackLength int    //Temporary for convenience
	Repos       string `yaml:"rest"`
}

type WorkoutVMA struct {
	VMA       string
	TotalTime string
	TimeTrack string
	Speed     string
	Pace      string
}

type Program struct {
	Name     string
	Date     time.Time
	Comment  string
	Workouts []Workout
}

type CalculatedProgram struct {
	Workout Workout
	VMAs    map[string]WorkoutVMA
}

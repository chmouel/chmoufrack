package main

import "time"

var (
	TRACK_LENGTH = 400
	VMA          = []int{13, 17}
	CONFIG_DIR   = ""
)

type Workout struct {
	ID          int64
	Repetition  string
	Meters      string
	Percentage  string
	TrackLaps   string
	TrackLength int
	Repos       string
}

type WorkoutVMA struct {
	VMA       string
	TotalTime string
	TimeTrack string
	Speed     string
	Pace      string
}

type Program struct {
	ID      int64
	Name    string
	Date    time.Time
	Comment string
}

type TemplateStruct struct {
	ProgramName string
	WP          Workout
	VMAs        map[string]WorkoutVMA
}

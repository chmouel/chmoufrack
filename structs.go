package main

import "time"

var (
	TRACK_LENGTH = 400
	VMA          = []int{13, 17}
)

type Workout struct {
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
	Name    string
	Date    time.Time
	Comment string
}

type TemplateStruct struct {
	WP   Workout
	VMAs map[string]WorkoutVMA
}

package main

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

type TemplateStruct struct {
	WP   Workout
	VMAs map[string]WorkoutVMA
}

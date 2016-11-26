package main

import (
	"database/sql"
	"time"
)

var (
	TRACK_LENGTH = 400
	VMA          = []int{14, 19}
	CONFIG_DIR   = ""
	DB           *sql.DB
)

type Workout struct {
	ID          int64
	Repetition  string `yaml:"laps"`
	Meters      string `yaml:"length"`
	Percentage  string `yaml:"percentage"`
	TrackLaps   string
	TrackLength int
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

package chmoufrack

import (
	"database/sql"
	"log"
	"os/user"
	"path/filepath"
	"time"
)

// getConfigDir ...
func getConfigDir() (configDir string) {
	user, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return filepath.Join(user.HomeDir, ".config", "frack")
}

var (
	TrackLength = 400
	TargetVma   = "14:19"
	ConfigDir   = getConfigDir()
	DB          *sql.DB
	StaticDir   = filepath.Join(ConfigDir, "static")
)

type Workout struct {
	ID          int64
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
	ID      int64
	Name    string
	Date    time.Time
	Comment string
}

type TemplateStruct struct {
	ProgramName string
	Workout     Workout
	VMAs        map[string]WorkoutVMA
}

package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"text/template"
)

func generate_content(ts TemplateStruct, content *bytes.Buffer) (err error) {
	fmt.Println(ts.WP.Repetition)
	t, err := template.ParseFiles("templates/content.tmpl")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = t.Execute(content, ts)
	if err != nil {
		return
	}
	return
}

func generate_template(content string) (err error) {
	dico := map[string]string{
		"Content": content,
	}

	t, err := template.ParseFiles("templates/template.tmpl")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = t.Execute(os.Stdout, dico)
	if err != nil {
		return
	}
	return
}

func generate_html(rounds [][]string) error {
	var content bytes.Buffer

	for i := range rounds {
		w := Workout{}
		w.TrackLength = TRACK_LENGTH

		repetition, err := strconv.ParseFloat(rounds[i][0], 64)
		if err != nil {
			return err
		}
		w.Repetition = fmt.Sprintf("%.f", repetition)

		percentage, err := strconv.ParseFloat(rounds[i][1], 64)
		if err != nil {
			return err
		}
		w.Percentage = fmt.Sprintf("%.f", percentage)
		meters, err := strconv.ParseFloat(rounds[i][2], 64)
		if err != nil {
			log.Fatal(err)
		}
		w.Meters = fmt.Sprintf("%.f", meters)

		w.Repos = rounds[i][3]

		track_length, err := strconv.ParseFloat(strconv.Itoa(TRACK_LENGTH), 64)
		if err != nil {
			return err
		}
		track_laps := meters / track_length
		laps := fmt.Sprintf("%.1f", track_laps)
		if strings.HasSuffix(laps, ".0") {
			laps = strings.TrimSuffix(laps, ".0")
		} else if laps == "0.5" {
			laps = "½"
		} else if strings.HasSuffix(laps, ".5") {
			laps = strings.Replace(laps, ".5", "½", -1)
		}
		w.TrackLaps = laps

		vmas := map[string]WorkoutVMA{}
		for _, vmad := range VMA {
			wt := WorkoutVMA{}
			vma := float64(vmad)
			total_time, err := calcul_vma_distance(vma, percentage, meters)
			if err != nil {
				return err
			}
			wt.VMA = fmt.Sprintf("%.f", vma)
			wt.TotalTime = total_time
			if int(meters) >= TRACK_LENGTH {
				time_lap, err := calcul_vma_distance(vma, percentage, float64(TRACK_LENGTH))
				if err != nil {
					return err
				}
				wt.TimeTrack = time_lap
			} else {
				wt.TimeTrack = "NA"
			}
			speed := calcul_vma_speed(vma, percentage)
			wt.Speed = fmt.Sprintf("%.2f", speed)
			wt.Pace = calcul_pace(speed)

			vmas[wt.VMA] = wt
		}
		ts := TemplateStruct{
			VMAs: vmas,
			WP:   w,
		}

		err = generate_content(ts, &content)
		if err != nil {
			return err
		}
		//spew.Dump(w)
	}
	err := generate_template(content.String())
	if err != nil {
		return err
	}
	return nil
}

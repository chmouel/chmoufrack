package chmoufrack

import (
	"fmt"
	"math"
	"strconv"
)

// round float64 - cause the stdlib doesnt have it :(
func round(val float64) float64 {
	var newVal float64

	if math.Abs(val-math.Ceil(val)) <= 0.5 {
		newVal = math.Ceil(val)
	} else {
		newVal = math.Floor(val)
	}
	return newVal
}

// calcul vma of a distance, you give a vma and a percent for a distance
func calcul_vma_distance(vma, percent, distance float64) (result string, err error) {
	vma_ms := vma * 1000 / 3600
	vma_100 := 100 / vma_ms
	calcul := vma_100 / percent * distance

	stemps := int(calcul)
	minute := ((stemps - (stemps)%60) / 60)
	second := ((stemps % 60) * 10) / 10

	if minute > 0 {
		result += strconv.Itoa(minute) + "'"
	}
	result += fmt.Sprintf("%.2d", second)
	if minute == 0 {
		result += "s"
	}

	return
}

// calcul_vma_vitesse ...
func calcul_vma_speed(vma, percent float64) float64 {
	return (vma * percent) / 100
}

// calcul_pace from a speed (kmh)
func calcul_pace(vitesse float64) (ret string) {
	var e = 1 / vitesse * 60
	var t = math.Floor(e / 60)
	var n = math.Floor(e - t*60)
	var r = round(60 * (e - t*60 - n))

	if r == 60 {
		n += 1
		r = 0
	}

	if n == 0 && r != 0 {
		return fmt.Sprintf("%0.f\"", r)
	}

	ret += fmt.Sprintf("%0.f'", n)
	if r == 0 {
		return
	} else if r < 10 {
		ret += "0"
	}

	ret += fmt.Sprintf("%0.f", r)

	return
}

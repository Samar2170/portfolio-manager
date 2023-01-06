package utils

import (
	"time"
)

func GetNextQuarter(t time.Time) time.Time {
	y, m, _ := t.Date()
	qm, qy := findNearestQNext(y, int(m))
	nd := time.Date(qy, time.Month(qm), 1, 0, 0, 0, 0, time.UTC)
	return nd
}

func GetNextHY(t time.Time) time.Time {
	var hm, hy int
	y, m, _ := t.Date()
	if m <= 6 {
		hy, hm = y, 7
	} else {
		hy, hm = y+1, 1
	}
	nd := time.Date(hy, time.Month(hm), 1, 0, 0, 0, 0, time.UTC)
	return nd
}

func findNearestQNext(y int, m int) (int, int) {
	switch m {
	case 1, 2, 3:
		return 4, y
	case 4, 5, 6:
		return 7, y
	case 7, 8, 9:
		return 10, y
	case 10, 11, 12:
		return 1, y + 1
	default:
		return 0, 0
	}

}

package utils

import (
	"errors"
	"math"
	"regexp"
	"strconv"
	"time"
)

// TimeAgo returns the rounded form of amount of time elapsed between now and the passed time
func TimeAgo(passedTime time.Time) string {

	// Get the elapsed number of hours
	floatHoursAgo := time.Since(passedTime).Seconds()

	// Take abs value since we handle the "ago" part below
	floatHoursAgo = math.Abs(floatHoursAgo)

	// Round the float
	secondsAgo := Round(floatHoursAgo)

	switch {
	case secondsAgo < 60:
		return strconv.Itoa(secondsAgo) + " seconds ago"
	case secondsAgo >= 60 && secondsAgo < 3600:
		return strconv.Itoa(secondsAgo/60) + " minutes ago"
	case secondsAgo >= 3600 && secondsAgo < 86400:
		if secondsAgo/3600 == 1 {
			return "1 hour ago"
		}
		return strconv.Itoa(secondsAgo/3600) + " hours ago"
	case secondsAgo >= 86400:
		if secondsAgo/86400 == 1 {
			return "1 day ago"
		}
		return strconv.Itoa(secondsAgo/86400) + " days ago"
	}

	return ""
}

// Round rounds the float to the nearest int
func Round(f float64) int {
	if math.Abs(f) < 0.5 {
		return 0
	}
	return int(f + math.Copysign(0.5, f))
}

// StatToSeconds takes in a beego stat param (e.g 20.40us or 15.20ms) and returns the time in seconds
func StatToSeconds(stat string) (float64, error) {

	// First parse out the ms, s, us, and the amount
	r := regexp.MustCompile("([0-9]+.[0-9]+)([a-z]+)")

	results := r.FindStringSubmatch(stat)

	if len(results) > 1 {
		valueStr := results[1]
		value, err := strconv.ParseFloat(valueStr, 10)
		if err != nil {
			return 0, err
		}

		time := results[2]

		switch time {
		case "us":
			convValue := value / 1000000
			// microseconds to seconds
			return convValue, nil
		case "ms":
			convValue := value / 1000
			// microseconds to seconds
			return convValue, nil
			// milliseconds to seconds
		case "s":
			return value, nil
		}
	}

	return 0, errors.New("Failed to parse time string from beego")

}

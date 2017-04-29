package settings

import (
	"encoding/json"
	"errors"
	"fmt"
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
	r := regexp.MustCompile("([0-9]+[.]?[0-9]*)([a-z]{1,2})")

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
			// microseconds to seconds
			convValue := value / 1000000
			return convValue, nil
		case "ms":
			// milliseconds to seconds
			convValue := value / 1000
			return convValue, nil
		case "s":
			return value, nil
		case "m":
			// minutes to seconds
			convValue := value * 60
			return convValue, nil
		case "h":
			// hours to seconds
			convValue := value * 3600
			return convValue, nil
		}
	}

	return 0, errors.New("Failed to parse time string from beego")

}

func Dump(obj interface{}) {
	b, _ := json.MarshalIndent(obj, "", "  ")
	fmt.Println(string(b))
}

package utils

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

// TestDaysAgo tests the DaysAgo function to check that the rounded timeAgo is equal to what we expecte
func TestDaysAgo(t *testing.T) {

	// Test days
	timeAgo := TimeAgo(time.Now().AddDate(0, 0, 4))
	Convey("TimeAgo should return the number of days ago", t, func() {
		So(timeAgo, ShouldEqual, "4 days ago")
	})

	// Test seconds
	timeAgo = TimeAgo(time.Now().Add(-50 * time.Second))
	Convey("TimeAgo should return the number of seconds ago", t, func() {
		So(timeAgo, ShouldEqual, "50 seconds ago")
	})

	// Test hours
	timeAgo = TimeAgo(time.Now().Add(-100 * time.Minute))
	Convey("TimeAgo should return the number of hours ago", t, func() {
		So(timeAgo, ShouldEqual, "1 hour ago")
	})

	// Test minutes
	timeAgo = TimeAgo(time.Now().Add(-10 * time.Minute))
	Convey("TimeAgo should return the number of minutes ago", t, func() {
		So(timeAgo, ShouldEqual, "10 minutes ago")
	})

	// Test minutes
	timeAgo = TimeAgo(time.Now().Add(-82001 * time.Second))
	Convey("TimeAgo should return the number of hours ago", t, func() {
		So(timeAgo, ShouldEqual, "22 hours ago")
	})
}

// TestRound tests the the simple round function
func TestRound(t *testing.T) {

	result := Round(3.99999)
	Convey("Round should return the expected result for 3.99999", t, func() {
		So(result, ShouldEqual, 4)
	})

	result = Round(0.4)
	Convey("Round should return the expected result for 0.4", t, func() {
		So(result, ShouldEqual, 0)
	})

}

// TestStatToSeconds tests the StatToSeconds function for the expected seconds return
func TestStatToSeconds(t *testing.T) {

	// Test minutes
	seconds, err := StatToSeconds("10.25ms")
	Convey("StatToSeconds should return the expected seconds result for 10ms", t, func() {
		So(seconds, ShouldEqual, 0.01025)
		So(err, ShouldBeNil)
	})

	seconds, err = StatToSeconds("20.39us")
	Convey("StatToSeconds should return the expected seconds result for 20.39us", t, func() {
		So(seconds, ShouldEqual, 2.039e-05)
		So(err, ShouldBeNil)
	})

	seconds, err = StatToSeconds("400s")
	Convey("StatToSeconds should return the expected seconds result for 400s", t, func() {
		So(seconds, ShouldEqual, 400)
		So(err, ShouldBeNil)
	})

	seconds, err = StatToSeconds("4000000seconds")
	Convey("StatToSeconds should return an error for an invalid result", t, func() {
		So(seconds, ShouldEqual, 0)
		So(err, ShouldNotBeNil)
	})
}

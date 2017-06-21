package funcs

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

// TestDaysAgo tests the DaysAgo function to check that the rounded timeAgo is equal to what we expecte
func TestDaysAgo(t *testing.T) {

	// Test days
	timeAgo := TimeAgo(time.Now().AddDate(0, 0, 4))
	Convey("TimeAgo should return the number of days ago", t, func(c C) {
		c.So(timeAgo, ShouldEqual, "4 days ago")
	})

	// Test seconds
	timeAgo = TimeAgo(time.Now().Add(-50 * time.Second))
	Convey("TimeAgo should return the number of seconds ago", t, func(c C) {
		c.So(timeAgo, ShouldEqual, "50 seconds ago")
	})

	// Test hours
	timeAgo = TimeAgo(time.Now().Add(-100 * time.Minute))
	Convey("TimeAgo should return the number of hours ago", t, func(c C) {
		c.So(timeAgo, ShouldEqual, "1 hour ago")
	})

	// Test minutes
	timeAgo = TimeAgo(time.Now().Add(-10 * time.Minute))
	Convey("TimeAgo should return the number of minutes ago", t, func(c C) {
		c.So(timeAgo, ShouldEqual, "10 minutes ago")
	})

	// Test minutes
	timeAgo = TimeAgo(time.Now().Add(-82001 * time.Second))
	Convey("TimeAgo should return the number of hours ago", t, func(c C) {
		c.So(timeAgo, ShouldEqual, "22 hours ago")
	})
}

// TestRound tests the the simple round function
func TestRound(t *testing.T) {

	result := Round(3.99999)
	Convey("Round should return the expected result for 3.99999", t, func(c C) {
		c.So(result, ShouldEqual, 4)
	})

	result = Round(0.4)
	Convey("Round should return the expected result for 0.4", t, func(c C) {
		c.So(result, ShouldEqual, 0)
	})

}

// TestStatToSeconds tests the StatToSeconds function for the expected seconds return
func TestStatToSeconds(t *testing.T) {

	// Test minutes
	seconds, err := StatToSeconds("10.25ms")
	Convey("StatToSeconds should return the expected seconds result for 10ms", t, func(c C) {
		c.So(seconds, ShouldEqual, 0.01025)
		c.So(err, ShouldBeNil)
	})

	seconds, err = StatToSeconds("20.39us")
	Convey("StatToSeconds should return the expected seconds result for 20.39us", t, func(c C) {
		c.So(seconds, ShouldEqual, 2.039e-05)
		c.So(err, ShouldBeNil)
	})

	seconds, err = StatToSeconds("4s")
	Convey("StatToSeconds should return the expected seconds result for 4s", t, func(c C) {
		c.So(seconds, ShouldEqual, 4)
		c.So(err, ShouldBeNil)
	})

	seconds, err = StatToSeconds("11m")
	Convey("StatToSeconds should return the expected seconds result for 11m", t, func(c C) {
		c.So(seconds, ShouldEqual, 660)
		c.So(err, ShouldBeNil)
	})

	seconds, err = StatToSeconds("100h")
	Convey("StatToSeconds should return the expected seconds result for 100h", t, func(c C) {
		c.So(seconds, ShouldEqual, 360000)
		c.So(err, ShouldBeNil)
	})

	seconds, err = StatToSeconds("4000000seconds")
	Convey("StatToSeconds should return an error for an invalid result", t, func(c C) {
		c.So(seconds, ShouldEqual, 0)
		c.So(err, ShouldNotBeNil)
	})
}

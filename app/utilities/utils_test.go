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

	// Test minutes
	timeAgo = TimeAgo(time.Now().Add(-100 * time.Minute))
	Convey("TimeAgo should return the number of hours ago", t, func() {
		So(timeAgo, ShouldEqual, "1 hour ago")
	})
}

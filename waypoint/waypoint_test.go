// package waypoint todo
package waypoint

import (
	"testing"

	"go.viam.com/test"
)


func TestReadWaypointFile(t *testing.T){
    waypoints, err := ReadWaypointsFromFile("./testwaypoints.csv")
    test.That(t, err, test.ShouldBeNil)
    test.That(t, waypoints , test.ShouldHaveLength, 3)
    test.That(t, waypoints[0].X, test.ShouldEqual, 0)
    test.That(t, waypoints[0].Y, test.ShouldEqual, 1)
    test.That(t, waypoints[0].Z, test.ShouldEqual, 2)
    test.That(t, waypoints[2].X, test.ShouldAlmostEqual, -0.543)
}


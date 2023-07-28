// Package waypoint contains code for represting waypoints
package waypoint

import (
	"os"

	"github.com/gocarina/gocsv"
)

// Waypoint represents a point on a slam map we will navigate to.
type Waypoint struct {
	X float64 `csv:"x"`
	Y float64 `csv:"y"`
	Z float32 `csv:"z"`
}

// ReadWaypointsFromFile reads a csv and turns it into a list of points.
func ReadWaypointsFromFile(path string) ([]*Waypoint, error) {
	//nolint:gosec
	waypointFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	//nolint:errcheck,gosec
	defer waypointFile.Close()
	waypoints := []*Waypoint{}

	if err := gocsv.UnmarshalFile(waypointFile, &waypoints); err != nil {
		return nil, err
	}
	return waypoints, nil
}

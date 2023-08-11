// Package waypoint contains code for represting waypoints
package waypoint

import (
	"os"

	"github.com/gocarina/gocsv"
	"github.com/golang/geo/r3"
)

// Waypoint represents a point on a slam map we will navigate to.
type Waypoint struct {
	X float64 `csv:"x"`
	Y float64 `csv:"y"`
	Z float64 `csv:"z"`
}

func (wp *Waypoint) ConvertToR3Vector() r3.Vector {
	// coordinates from web ui are in m, MoveOnMap() uses mm so need to convert
	return r3.Vector{
		X: wp.X * 1000.0,
		Y: wp.Y * 1000.0,
		Z: 0,
	}
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

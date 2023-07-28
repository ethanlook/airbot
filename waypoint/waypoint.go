// package waypoint todo
package waypoint

import (
	"github.com/gocarina/gocsv"
	"os"
)

type Waypoint struct {
	X float64 `csv:"x"`
	Y float64 `csv:"y"`
	Z float32 `csv:"z"`
}

func ReadWaypointsFromFile(path string) ([]*Waypoint, error) {
	waypointFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	//nolint:errcheck
	defer waypointFile.Close()
	waypoints := []*Waypoint{}

	if err := gocsv.UnmarshalFile(waypointFile, &waypoints); err != nil {
		return nil, err
	}
	return waypoints, nil
}

// Package airbot runs the main navigation loop and data collection
package airbot

import (
	"errors"
	"fmt"

	"github.com/edaniels/golog"
	"github.com/ethanlook/airbot/imagedetector"
	"github.com/ethanlook/airbot/move"
	pb "github.com/ethanlook/airbot/proto/v1"
	"github.com/ethanlook/airbot/waypoint"

	"go.viam.com/rdk/robot/client"
)

var errRouteUnspecified = errors.New("route unspecified")

// AirBot is the main navigation loop and data collection.
type AirBot struct {
	logger      golog.Logger
	robotClient *client.RobotClient
}

// NewAirBot creates a new AirBot.
func NewAirBot(logger golog.Logger, robotClient *client.RobotClient) *AirBot {
	return &AirBot{
		logger,
		robotClient,
	}
}

// Start starts the main navigation loop and data collection.
func (a *AirBot) Start(route pb.Route) error {
	var waypointsFile string
	switch route {
	case pb.Route_ROUTE_KITCHEN:
		waypointsFile = "./routes/kitchen-route.csv"
	case pb.Route_ROUTE_UNSPECIFIED:
		fallthrough
	default:
		return errRouteUnspecified
	}
	waypoints, err := waypoint.ReadWaypointsFromFile(waypointsFile)
	if err != nil {
		return fmt.Errorf("error reading waypoints from file: %w", err)
	}

	moveManager, err := move.NewMoveManager(a.robotClient, a.logger)
	if err != nil {
		return fmt.Errorf("error creating move manager: %w", err)
	}
	detector, err := imagedetector.NewDetector(a.robotClient, "top-cam", "coffee-mug-detector", a.logger)
	if err != nil {
		return fmt.Errorf("error creating detector: %w", err)
	}

	for i, w := range waypoints {
		a.logger.Infof("Starting navigation to waypoint #%d: %w", i, w)
		err := moveManager.MoveOnMap(w, 3)
		if err != nil {
			return fmt.Errorf("error moving on map: %w", err)
		}

		a.logger.Infof("Successfully made it to waypoint: %w", w)

		a.logger.Info("Starting coffee mug detection")

		for j := 1; j <= 4; j++ {
			a.logger.Infof("Turning 90 degrees #%d", j)
			err = moveManager.Turn90()
			if err != nil {
				a.logger.Errorw("error turning 90 degrees", "err", err)
			}

			a.logger.Info("Doing image detection")
			detections, err := detector.GetDetectionsFromCamera()
			if err != nil {
				continue
			}

			a.logger.Infof("Found %d mugs at waypoint #%d, turn #%d", detector.HowManyMugs(detections), i, j)
		}

		a.logger.Info("Finished coffee mug detection")
	}

	return nil
}

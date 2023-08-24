// Package airbot runs the main navigation loop and data collection
package airbot

import (
	"github.com/edaniels/golog"
	"github.com/ethanlook/airbot/imagedetector"
	"github.com/ethanlook/airbot/move"
	"github.com/ethanlook/airbot/waypoint"

	"go.viam.com/rdk/robot/client"
)

// AirBot is the main navigation loop and data collection.
type AirBot struct {
	logger      golog.Logger
	robotClient *client.RobotClient
	waypoints   []*waypoint.Waypoint
}

// NewAirBot creates a new AirBot.
func NewAirBot(logger golog.Logger, robotClient *client.RobotClient, waypoints []*waypoint.Waypoint) *AirBot {
	return &AirBot{
		logger:      logger,
		robotClient: robotClient,
		waypoints:   waypoints,
	}
}

// Start starts the main navigation loop and data collection.
func (a *AirBot) Start() {
	moveManager, err := move.NewMoveManager(a.robotClient, a.logger)
	if err != nil {
		a.logger.Errorw("error creating move manager", "err", err)
		return
	}
	detector, err := imagedetector.NewDetector(a.robotClient, "top-cam", "coffee-mug-detector", a.logger)
	if err != nil {
		a.logger.Errorw("error creating detector", "err", err)
		return
	}

	for i, w := range a.waypoints {
		a.logger.Infof("Starting navigation to waypoint #%d: %w", i, w)
		err := moveManager.MoveOnMap(w, 3)
		if err != nil {
			a.logger.Errorw("error moving on map", "err", err)
			a.logger.Errorw("exiting the program")
			return
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
}

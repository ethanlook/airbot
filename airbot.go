// Package airbot runs the main navigation loop and data collection
package airbot

import (
	"github.com/edaniels/golog"
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
<<<<<<< HEAD
	// motion.FromRobot()
	// slam, err := slam.FromRobot(a.robotClient, "slam-service")
	for _, w := range a.waypoints {
		a.logger.Infof("Navigating to waypoint: %v", w)
		// slam.
=======
	for _, w := range a.waypoints {
		a.logger.Infof("Waypoint: %v", w)
>>>>>>> main
	}
}

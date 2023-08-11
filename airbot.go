// Package airbot runs the main navigation loop and data collection
package airbot

import (
	"context"
	"fmt"
	"math"

	"github.com/edaniels/golog"
	"github.com/ethanlook/airbot/waypoint"
	"github.com/pkg/errors"

	"go.viam.com/rdk/components/base"
	"go.viam.com/rdk/robot/client"
	"go.viam.com/rdk/services/slam"
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
	// motion.FromRobot()
	// slam, err := slam.FromRobot(a.robotClient, "slam-service")
	ctx := context.Background()
	for _, w := range a.waypoints {
		a.logger.Infof("Navigating to waypoint: %v", w)
		// slam.
	}
	a.GetPos(ctx)
	err := a.tryMoveToWaypoint(ctx, waypoint.Waypoint{
		X: 0.26,
		Y: 3.25,
	}, 0.1, 1)
	fmt.Println(err)
}

func (a *AirBot) GetPos(ctx context.Context) (*waypoint.Waypoint, float64, error) {
	slam, err := slam.FromRobot(a.robotClient, "slam-service")
	if err != nil {
		return nil, 0.0, err
	}
	pos, _, err := slam.GetPosition(ctx)
	fmt.Println(pos)
	point := pos.Point()
	fmt.Printf("pos: %v\n", pos.Point())
	fmt.Printf("or: %v\n", pos.Orientation().AxisAngles().Theta)
	return &waypoint.Waypoint{
		X: point.X / 1000.0,
		Y: point.Y / 1000.0,
	}, pos.Orientation().AxisAngles().Theta, nil
}

func (a *AirBot) distAndAngleTo(ctx context.Context, desiredPos waypoint.Waypoint) (float64, float64, error) {
	currentPos, theta, err := a.GetPos(ctx)
	if err != nil {
		return 0, 0, err
	}
	dx := desiredPos.X - currentPos.X
	dy := desiredPos.Y - currentPos.Y
	distSquared := dx*dx + dy*dy
	desiredTheta := math.Atan2(dy, dx)
	detaTheta := desiredTheta - theta
	dist := math.Sqrt(distSquared)
	return dist, detaTheta, nil
}

func (a *AirBot) tryMoveToWaypoint(ctx context.Context, desiredPos waypoint.Waypoint, tol float64, numTries int) error {
	n := 0
	base, err := base.FromRobot(a.robotClient, "viam_base")
	if err != nil {
		return err
	}
	dist, detaTheta, err := a.distAndAngleTo(ctx, desiredPos)
	if err != nil {
		return err

	}
	for dist >= tol && n < numTries {
		n += 1

		fmt.Printf("starting spin %v degrees\n", int(detaTheta*57.29))
		err = base.Spin(ctx, detaTheta*57.29, 100, map[string]interface{}{})
		if err != nil {
			return err
		}
		fmt.Printf(" starting move %v\n", dist)
		err = base.MoveStraight(ctx, int(math.Floor(dist*1000)), 1000.0, map[string]interface{}{})
		if err != nil {
			return err
		}
		dist, detaTheta, err = a.distAndAngleTo(ctx, desiredPos)
		if err != nil {
			return err

		}
	}
	if dist <= tol {
		return nil

	}
	return errors.New("failed to reach desired tol")

}

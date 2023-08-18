package move

import (
	"context"

	"github.com/edaniels/golog"
	"github.com/ethanlook/airbot/waypoint"
	"go.viam.com/rdk/components/base"
	"go.viam.com/rdk/robot/client"
	"go.viam.com/rdk/services/motion"
	"go.viam.com/rdk/services/slam"
	"go.viam.com/rdk/spatialmath"
)

type Move interface {
	MoveOnMap(wp *waypoint.Waypoint) error
}

type MoveManger struct {
	// robot client
	rc *client.RobotClient
	// motion service
	ms motion.Service
	// slam service
	slam slam.Service
	// allows us to cancel the request
	base   base.Base
	logger golog.Logger
}

func NewMoveManager(robotClient *client.RobotClient, logger golog.Logger) (Move, error) {
	ms, err := motion.FromRobot(robotClient, "builtin")
	if err != nil {
		return nil, err
	}
	slam, err := slam.FromRobot(robotClient, "slam-kitchen3")
	if err != nil {
		return nil, err
	}
	base, err := base.FromRobot(robotClient, "viam_base")
	if err != nil {
		return nil, err
	}

	return &MoveManger{rc: robotClient, ms: ms, slam: slam, base: base, logger: logger}, nil
}

func (mm *MoveManger) MoveOnMap(wp *waypoint.Waypoint) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// moveOnMap uses the last orientation to pass to the new pose
	// so followed that implementation there

	lastPose, _, err := mm.slam.GetPosition(ctx)
	if err != nil {
		return err
	}

	pose := spatialmath.NewPose(wp.ConvertToR3Vector(), lastPose.Orientation())
	motionConfig := make(map[string]interface{})
	motionConfig["motion_profile"] = "position_only"
	motionConfig["timeout"] = 5

	_, err = mm.ms.MoveOnMap(ctx, mm.base.Name(), pose, mm.slam.Name(), motionConfig)
	if err != nil {
		return err
	}
	return nil
}

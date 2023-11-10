// Package move contains the code to make a robot move on a slam map
package move

import (
	"context"

	"github.com/ethanlook/airbot/waypoint"

	"go.viam.com/rdk/components/base"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/resource"
	"go.viam.com/rdk/services/motion"
	"go.viam.com/rdk/services/slam"
	"go.viam.com/rdk/spatialmath"
)

// Move defines the interface to Move.
type Move interface {
	MoveOnMap(wp *waypoint.Waypoint, attempts int) error
	Turn90() error
}

// Manager holds all necessary info to move.
type Manager struct {
	// motion service
	ms motion.Service
	// slam service
	slam slam.Service
	// allows us to cancel the request
	base   base.Base
	logger logging.Logger
}

// NewMoveManager creates a MoveManager.
func NewMoveManager(logger logging.Logger, deps resource.Dependencies, slamService string, baseComponent string) (Move, error) {
	ms, err := motion.FromDependencies(deps, "builtin")
	if err != nil {
		return nil, err
	}
	slam, err := resource.FromDependencies[slam.Service](deps, resource.NewName(slam.API, slamService))
	if err != nil {
		return nil, err
	}
	base, err := base.FromDependencies(deps, baseComponent)
	if err != nil {
		return nil, err
	}

	return &Manager{ms: ms, slam: slam, base: base, logger: logger}, nil
}

// MoveOnMap moves the rover to a waypoint on the slam map.
func (mm *Manager) MoveOnMap(wp *waypoint.Waypoint, attempts int) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for i := 1; i <= attempts; i++ {
		// moveOnMap uses the last orientation to pass to the new pose
		// so followed that implementation there

		// TODO(ethanlook): Child context?

		lastPose, _, err := mm.slam.Position(ctx)
		if err != nil {
			return err
		}

		pose := spatialmath.NewPose(wp.ConvertToR3Vector(), lastPose.Orientation())
		motionConfig := make(map[string]interface{})
		motionConfig["motion_profile"] = "position_only"
		motionConfig["timeout"] = 30
		_, err = mm.ms.MoveOnMap(ctx, mm.base.Name(), pose, mm.slam.Name(), motionConfig)
		if err != nil && i == attempts {
			mm.logger.Errorw("Errored on final attempt to navigate to waypoint", "err", err)
			return err
		} else if err == nil {
			return nil
		}
		mm.logger.Errorw("Navigation attempt failed", "err", err)
		mm.logger.Infof("Retry navigation to waypoint (#%d): %w", i, wp)
	}
	return nil
}

// Turn90 spins the base 90 degrees.
func (mm *Manager) Turn90() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	return mm.base.Spin(ctx, 90, 45, map[string]interface{}{})
}

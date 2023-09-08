// Package airbot runs the main navigation loop and data collection
package airbot

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/edaniels/golog"
	"github.com/ethanlook/airbot/imagedetector"
	"github.com/ethanlook/airbot/move"
	pb "github.com/ethanlook/airbot/proto/v1"
	"github.com/ethanlook/airbot/waypoint"
	"github.com/pkg/errors"

	"go.viam.com/rdk/components/generic"
	"go.viam.com/rdk/resource"
)

var errRouteUnspecified = errors.New("route unspecified")

var Model = resource.NewModel("ethanlook", "service", "airbot")

func init() {
	registration := resource.Registration[resource.Resource, *Config]{
		Constructor: func(ctx context.Context, deps resource.Dependencies, conf resource.Config, logger golog.Logger) (resource.Resource, error) {
			return newAirBot(ctx, deps, conf, logger)
		},
	}
	resource.RegisterComponent(generic.API, Model, registration)
}

func newAirBot(ctx context.Context, deps resource.Dependencies, conf resource.Config, logger golog.Logger) (resource.Resource, error) {
	newConf, err := resource.NativeConfig[*Config](conf)
	if err != nil {
		return nil, errors.Wrap(err, "create component failed due to config parsing")
	}

	moveManager, err := move.NewMoveManager(logger, deps, newConf.SlamService, newConf.BaseComponent)
	if err != nil {
		return nil, fmt.Errorf("error creating move manager: %w", err)
	}
	detector, err := imagedetector.NewDetector(logger, deps, newConf.VisionService, newConf.CameraComponent)
	if err != nil {
		return nil, fmt.Errorf("error creating image detector: %w", err)
	}

	cancelCtx, cancelFunc := context.WithCancel(context.Background())

	instance := &AirBot{
		Named:  conf.ResourceName().AsNamed(),
		config: newConf,
		deps:   deps,

		cancelCtx:  cancelCtx,
		cancelFunc: cancelFunc,

		logger: logger,

		isRunning:    false,
		isRunningMux: sync.RWMutex{},

		moveManager: moveManager,
		detector:    detector,
	}
	instance.logger.Infoln("Started")
	return instance, nil
}

// AirBot is the main navigation loop and data collection.
type AirBot struct {
	resource.Named
	resource.AlwaysRebuild
	config *Config
	deps   resource.Dependencies

	cancelCtx  context.Context
	cancelFunc func()

	logger golog.Logger

	isRunning       bool
	isRunningMux    sync.RWMutex
	startCancelFunc func()

	moveManager *move.MoveManager
	detector    *imagedetector.Detector
}

// DoCommand sends/receives arbitrary data
func (a *AirBot) DoCommand(ctx context.Context, cmd map[string]interface{}) (map[string]interface{}, error) {
	if startCmd, ok := cmd["start"]; ok {
		req, err := parseStartCommand(startCmd)
		if err != nil {
			return nil, err
		}
		return nil, a.Start(req)
	}

	if _, ok := cmd["stop"]; ok {
		return nil, a.Stop(&pb.StopRequest{})
	}

	return nil, fmt.Errorf("unknown DoCommand")
}

func parseStartCommand(cmd interface{}) (*pb.StartRequest, error) {
	start, startOk := cmd.(map[string]interface{})
	if !startOk {
		return nil, fmt.Errorf("could not cast start to map")
	}

	route, routeOk := start["route"]
	if !routeOk {
		return nil, fmt.Errorf("missing start.route")
	}

	routeStr, routeStrOk := route.(string)
	if !routeStrOk {
		return nil, fmt.Errorf("could not cast start.route to string")
	}

	startWaypointNum, startWaypointNumOk := start["start_waypoint_num"]
	if !startWaypointNumOk {
		startWaypointNum = uint32(0)
	}

	startWaypointNumInt, startWaypointNumIntOk := startWaypointNum.(uint32)
	if !startWaypointNumIntOk {
		return nil, fmt.Errorf("could not cast start.start_waypoint_num to uint32")
	}

	switch routeStr {
	case "kitchen":
		return &pb.StartRequest{
			Route:            pb.Route_ROUTE_KITCHEN,
			StartWaypointNum: startWaypointNumInt,
		}, nil

	case "test":
		return &pb.StartRequest{
			Route:            pb.Route_ROUTE_TEST,
			StartWaypointNum: startWaypointNumInt,
		}, nil
	}

	return nil, fmt.Errorf("unknown start.route")
}

// Close must safely shut down the resource and prevent further use.
// Close must be idempotent.
// Later reconfiguration may allow a resource to be "open" again.
func (a *AirBot) Close(ctx context.Context) error {
	a.logger.Info("close")
	a.cancelFunc()
	return nil
}

// Start starts the main navigation loop and data collection.
func (a *AirBot) Start(req *pb.StartRequest) error {
	if a.getIsRunning() {
		return fmt.Errorf("airbot is already running")
	}

	a.setIsRunning(true)

	ex, err := os.Executable()
	if err != nil {
		return err
	}

	var waypointsFile string
	switch req.Route {
	case pb.Route_ROUTE_KITCHEN:
		waypointsFile = "../../routes/kitchen-route.csv"
	case pb.Route_ROUTE_TEST:
		waypointsFile = "../../routes/test-route.csv"
	case pb.Route_ROUTE_UNSPECIFIED:
		fallthrough
	default:
		return errRouteUnspecified
	}

	waypoints, err := waypoint.ReadWaypointsFromFile(filepath.Join(ex, waypointsFile))
	if err != nil {
		return fmt.Errorf("error reading waypoints from file: %w", err)
	}

	if int(req.StartWaypointNum) > len(waypoints) {
		return fmt.Errorf("start_waypoint_num out of bounds, %d >= len(waypoints), len(waypoints) = %d", req.StartWaypointNum, len(waypoints))
	}

	a.startBgProcess(waypoints[req.StartWaypointNum:])

	return nil
}

func (a *AirBot) startBgProcess(waypoints []*waypoint.Waypoint) {
	ctx, cancel := context.WithCancel(a.cancelCtx)
	a.startCancelFunc = cancel

	go func(waypoints []*waypoint.Waypoint) {
		for i, w := range waypoints {
			a.logger.Infof("Starting navigation to waypoint #%d: %w", i, w)
			err := a.moveManager.MoveOnMap(ctx, w, 3)
			if err != nil {
				a.logger.Errorf("error moving on map: %w", err)
				a.setIsRunning(false)
				return
			}

			a.logger.Infof("Successfully made it to waypoint: %w", w)

			a.logger.Info("Starting coffee mug detection")

			for j := 1; j <= 4; j++ {
				a.logger.Infof("Turning 90 degrees #%d", j)
				err = a.moveManager.Turn90(ctx)
				if err != nil {
					a.logger.Errorw("error turning 90 degrees", "err", err)
				}

				a.logger.Info("Doing image detection")
				detections, err := a.detector.GetDetectionsFromCamera(ctx)
				if err != nil {
					continue
				}

				a.logger.Infof("Found %d mugs at waypoint #%d, turn #%d", a.detector.HowManyMugs(detections), i, j)
			}

			a.logger.Info("Finished coffee mug detection")
			a.setIsRunning(false)
		}
	}(waypoints)
}

func (a *AirBot) Stop(req *pb.StopRequest) error {
	if !a.getIsRunning() {
		return fmt.Errorf("airbot is not running")
	}
	a.startCancelFunc()
	a.setIsRunning(false)
	return nil
}

func (a *AirBot) getIsRunning() bool {
	a.isRunningMux.RLock()
	defer a.isRunningMux.RUnlock()
	return a.isRunning
}

func (a *AirBot) setIsRunning(isRunning bool) {
	a.isRunningMux.Lock()
	defer a.isRunningMux.Unlock()
	a.isRunning = isRunning
}

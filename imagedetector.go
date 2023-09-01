// Package imagedetector holds information regarding the detector
package airbot

import (
	"context"
	"fmt"

	"github.com/edaniels/golog"

	"go.viam.com/rdk/components/camera"
	"go.viam.com/rdk/resource"
	"go.viam.com/rdk/services/vision"
	"go.viam.com/rdk/vision/objectdetection"
)

// Detector stores information for the vision service detector.
type Detector struct {
	cam    camera.Camera
	vs     vision.Service
	logger golog.Logger
}

// GetDetectionsFromCamera gets detections from camera.
func (detector *Detector) GetDetectionsFromCamera() ([]objectdetection.Detection, error) {
	return detector.vs.DetectionsFromCamera(context.Background(), detector.cam.Name().Name, map[string]interface{}{})
}

// HowManyMugs determines how many mugs in the frame.
func (detector *Detector) HowManyMugs(detections []objectdetection.Detection) int {
	threshold := 0.7
	count := 0

	for _, detection := range detections {
		label := detection.Label()
		score := detection.Score()
		if label == "coffee-mug" && score > threshold {
			count++
		}
	}

	return count
}

// NewDetector returns a new detector.
func NewDetector(
	dependencies resource.Dependencies,
	airBot *AirBot,
	cameraName string,
	visionServiceName string,
	logger golog.Logger,
) (*Detector, error) {
	// Grab the camera from the robot
	myCam, err := camera.FromDependencies(dependencies, airBot.Config.CameraComponent)
	if err != nil {
		return nil, fmt.Errorf("cannot get camera: %w", err)
	}

	visService, err := vision.FromDependencies(dependencies, airBot.Config.VisionService)
	if err != nil {
		return nil, fmt.Errorf("cannot get vision service: %w", err)
	}

	return &Detector{
		cam:    myCam,
		vs:     visService,
		logger: logger,
	}, nil
}

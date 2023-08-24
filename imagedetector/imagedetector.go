package imagedetector

import (
	"context"

	"github.com/edaniels/golog"
	"go.viam.com/rdk/components/camera"
	"go.viam.com/rdk/robot/client"
	"go.viam.com/rdk/services/vision"
)

func NewDetector(robotClient *client.RobotClient, logger golog.Logger) error {
	// Grab the camera from the robot
	cameraName := "top-cam" // make sure to use the same component name that you have in your robot configuration
	myCam, err := camera.FromRobot(robotClient, cameraName)
	if err != nil {
		logger.Fatalf("cannot get camera: %v", err)
	}

	visService, err := vision.FromRobot(robotClient, "coffee-mug-vision-service")
	if err != nil {
		logger.Fatalf("Cannot get vision service: %v", err)
	}

	// Get detections from the camera output
	detections, err := visService.DetectionsFromCamera(context.Background(), cameraName, map[string]interface{}{})
	if err != nil {
		logger.Errorw("Could not get detections", "err", err)
		return err
	}
	if len(detections) > 0 {
		logger.Info(detections[0])
	}

	// If you need to store the image, get the image first
	// and then run detections on it. This process is slower:

	// Get the stream from a camera
	camStream, err := myCam.Stream(context.Background())
	if err != nil {
		logger.Errorw("Could not open camera stream", "err", err)
		return err
	}

	// Get an image from the camera stream
	img, release, err := camStream.Next(context.Background())
	if err != nil {
		logger.Errorw("Could not get an image from the camera stream", "err", err)
		return err
	}
	defer release()

	detectionsFromImage, err := visService.Detections(context.Background(), img, nil)
	if err != nil {
		logger.Errorw("Could not get detections", "err", err)
		return err
	}
	for _, detection := range detectionsFromImage {
		logger.Info(detection)
	}
	return nil
}

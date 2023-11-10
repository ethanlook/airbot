package airbot

import (
	"context"
	"time"

	"github.com/ethanlook/airbot/waypoint"
	datapb "go.viam.com/api/app/datasync/v1"
	"go.viam.com/rdk/resource"
	"go.viam.com/rdk/services/datamanager/datacapture"
	"go.viam.com/utils/protoutils"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// WaypointReading stores info about a waypoint.
type WaypointReading struct {
	// Position is the position of the waypoint to upload
	Position waypoint.Waypoint `json:"position"`
	// Pictures stores all of the pictures taken at the location
	Pictures []PhotoReading `json:"pictures"`

	// TempCelsius is the temp at the waypoint location
	TempCelsius float64 `json:"temp_c"`
	// Ppm25 is the ppm2.5 at the waypoint location
	Ppm25 float64 `json:"ppm25"`
}

// PhotoReading is a way to link a uploaded photo to a WaypointReading.
type PhotoReading struct {
	// This allows us to link a waypoint reading to an uploaded image.
	// It should be a UUID and should be passed into UploadPNG
	Tag string `json:"tag"`
	// MugsDetected stores the number of mugs detected in the image
	MugsDetected int64 `json:"mugs_detected"`
}

type AirbotDataManager struct {
	writer datacapture.BufferedWriter

	cancelCtx  context.Context
	cancelFunc context.CancelFunc
}

func NewAirbotDataManager(targetDir string) (*AirbotDataManager, error) {
	captureMetadata, err := datacapture.BuildCaptureMetadata(
		resource.NewAPI("airbot", "data", "collection"), //config.Name.API,
		"airbot",  //config.Name.ShortName(),
		"collect", //config.Method,
		nil,       //config.AdditionalParams,
		nil,       //config.Tags,
	)
	if err != nil {
		return nil, err

	}
	context, cancelFunc := context.WithCancel(context.Background())
	return &AirbotDataManager{
		cancelCtx:  context,
		cancelFunc: cancelFunc,

		writer: datacapture.NewBuffer(targetDir, captureMetadata),
	}, nil
}

// TODO
// Add a lock / basic async otherwise this may hang
func (adm *AirbotDataManager) UploadReading(wp *WaypointReading) error {
	proto, err := waypointReadingToProto(wp)
	if err != nil {
		return err
	}
	adm.writer.Write(proto)
	return nil
}

func (adm *AirbotDataManager) Close() error {
	adm.cancelFunc()
	return nil
}

func waypointReadingToProto(wp *WaypointReading) (*datapb.SensorData, error) {
	timeReceived := timestamppb.New(time.Now().UTC())

	pbReading, err := protoutils.StructToStructPb(wp)
	if err != nil {
		return nil, err
	}
	msg := datapb.SensorData{
		Metadata: &datapb.SensorMetadata{
			TimeRequested: timeReceived,
			TimeReceived:  timeReceived,
		},
		Data: &datapb.SensorData_Struct{
			Struct: pbReading,
		},
	}
	return &msg, nil
}

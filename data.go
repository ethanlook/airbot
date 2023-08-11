package airbot

import (
	"context"
	"time"

	"github.com/ethanlook/airbot/waypoint"
	v1 "go.viam.com/api/app/datasync/v1"
	"go.viam.com/rdk/resource"
	"go.viam.com/rdk/services/datamanager/datacapture"
	"go.viam.com/utils/protoutils"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type PhotoReading struct {
	Uri  string
	Info interface{}
}
type WaypointReading struct {
	ReadingName string
	Position    waypoint.Waypoint
	Pictures    []PhotoReading

	temp  float64
	ppm25 float64
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
	proto, err := WaypointReadingToProto(wp)
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

func WaypointReadingToProto(wp *WaypointReading) (*v1.SensorData, error) {
	timeRequested := timestamppb.New(time.Now().UTC())
	timeReceived := timestamppb.New(time.Now().UTC())

	pbReading, err := protoutils.StructToStructPb(wp)
	if err != nil {
		return nil, err
	}
	msg := v1.SensorData{
		Metadata: &v1.SensorMetadata{
			TimeRequested: timeRequested,
			TimeReceived:  timeReceived,
		},
		Data: &v1.SensorData_Struct{
			Struct: pbReading,
		},
	}
	return &msg, nil
}

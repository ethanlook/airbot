package airbot

import (
	"context"
	"io"
	"time"

	"github.com/ethanlook/airbot/waypoint"
	"github.com/pkg/errors"
	"go.uber.org/multierr"
	datapb "go.viam.com/api/app/datasync/v1"
	"go.viam.com/utils/protoutils"
	"go.viam.com/utils/rpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	// These are used for the upload metadata. To my knowledge, they can be pretty arbitrary.
	componentType = "Airbot"
	componentName = "Airbot"
	methodName    = "Exploring"
	// Used when uploading images.
	uploadChunkSize = 64 * 1024
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

// AirbotDataManager is a wrapper around the data services provided by viam.
//
//nolint:revive
type AirbotDataManager struct {
	cloudConn   *rpc.ClientConn
	robotPartID string
	dataClient  datapb.DataSyncServiceClient
	cancelCtx   context.Context
	cancelFunc  context.CancelFunc
}

// NewAirBotDataManager creates a new airbot data manager.
func NewAirBotDataManager(conn *rpc.ClientConn, robotPartID string) (*AirbotDataManager, error) {
	context, cancelFunc := context.WithCancel(context.Background())
	return &AirbotDataManager{
		cloudConn:   conn,
		robotPartID: robotPartID,
		dataClient:  datapb.NewDataSyncServiceClient(*conn),
		cancelCtx:   context,
		cancelFunc:  cancelFunc,
	}, nil
}

// UploadReading uploads a reading with tag of `routeName`.
func (dm *AirbotDataManager) UploadReading(wp *WaypointReading, routeName string) error {
	wpProto, err := waypointReadingToProto(wp)
	if err != nil {
		return err
	}

	uploadMeta := &datapb.UploadMetadata{
		PartId:        dm.robotPartID,
		ComponentType: componentType,
		ComponentName: componentName,
		MethodName:    methodName,
		Type:          datapb.DataType_DATA_TYPE_TABULAR_SENSOR,
		Tags:          []string{routeName},
	}

	uploadReq := &datapb.DataCaptureUploadRequest{
		Metadata:       uploadMeta,
		SensorContents: []*datapb.SensorData{wpProto},
	}
	_, err = dm.dataClient.DataCaptureUpload(context.Background(), uploadReq)
	return err
}

// Close closes the data manager.
func (dm *AirbotDataManager) Close() error {
	dm.cancelFunc()
	return nil
}

// UploadPNG expects bytes with a MIME type of image/png.
func (dm *AirbotDataManager) UploadPNG(photoName string, photo []byte, tag string) error {
	stream, err := dm.dataClient.FileUpload(dm.cancelCtx)
	if err != nil {
		return err
	}

	uploadMeta := &datapb.UploadMetadata{
		PartId:        dm.robotPartID,
		ComponentType: componentType,
		ComponentName: componentName,
		MethodName:    methodName,
		Type:          datapb.DataType_DATA_TYPE_FILE,
		FileName:      photoName,
		FileExtension: "png",
		Tags:          []string{tag},
	}

	req := &datapb.FileUploadRequest{
		UploadPacket: &datapb.FileUploadRequest_Metadata{
			Metadata: uploadMeta,
		},
	}
	if err := stream.Send(req); err != nil {
		return err
	}

	var errs error
	// We do not add the EOF as an error because all server-side errors trigger an EOF on the stream
	// This results in extra clutter to the error msg
	if err := dm.sendUploadRequests(stream, photo); err != nil && !errors.Is(err, io.EOF) {
		errs = multierr.Combine(errs, errors.Wrapf(err, "could not upload %s", photoName))
	}

	_, closeErr := stream.CloseAndRecv()
	return multierr.Combine(errs, closeErr)
}

func (dm *AirbotDataManager) sendUploadRequests(stream datapb.DataSyncService_FileUploadClient, bytes []byte) error {
	//nolint:errcheck
	defer stream.CloseSend()
	// Loop until there is no more content to be read from the array or the context expires.
	for i := 0; i < len(bytes); {
		if dm.cancelCtx.Err() != nil {
			return dm.cancelCtx.Err()
		}
		numBytesToSend := uploadChunkSize
		if i+numBytesToSend > len(bytes) {
			numBytesToSend = len(bytes) - i
		}
		next := bytes[i : i+numBytesToSend]
		i += numBytesToSend

		nextReq := &datapb.FileUploadRequest{
			UploadPacket: &datapb.FileUploadRequest_FileContents{
				FileContents: &datapb.FileData{
					Data: next,
				},
			},
		}

		if err := stream.Send(nextReq); err != nil {
			return err
		}
	}
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

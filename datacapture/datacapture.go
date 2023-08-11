// Package datacapture captures data
package datacapture

import (
	"go.viam.com/rdk/robot/client"
	"golang.org/x/net/context"
	datapb "go.viam.com/api/app/datasync/v1"
)

func UploadData(ctx context.Context, robotClient *client.RobotClient) error {
    client := datapb.DataSyncServiceClient
    datapb.NewDataSyncServiceClient(robotClient.Reconfigure)

    return nil
}

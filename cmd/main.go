// Package main starts the airbot.
package main

import (
	"context"
	"net/url"
	"os"

	"github.com/edaniels/golog"
	"github.com/ethanlook/airbot"
	"github.com/ethanlook/airbot/waypoint"
	"github.com/joho/godotenv"
	"go.viam.com/utils/rpc"

	"go.viam.com/rdk/robot/client"
	"go.viam.com/rdk/utils"
)
// var routeName = flag.String("route", "w1", "Route to follow")

func main() {
	logger := golog.NewDevelopmentLogger("client")
	ctx := context.Background()

	err := godotenv.Load()
	if err != nil {
		logger.Panic(err)
	}
	robotPartID := os.Getenv("ROBOT_PART_ID")
	cloudDialOpts := rpc.WithEntityCredentials(robotPartID, rpc.Credentials{
		Type:    utils.CredentialsTypeRobotSecret,
		Payload: os.Getenv("ROBOT_PART_SECRET"),
	})

	cloudURL, err := url.Parse(os.Getenv("CLOUD_ADDRESS"))
	if err != nil {
		logger.Fatal(err)
	}
	cloudConn, err := rpc.DialDirectGRPC(context.Background(), cloudURL.Host, logger, cloudDialOpts)
	if err != nil {
		logger.Fatal(err)
	}
	//nolint:errcheck
	defer cloudConn.Close()

	dataManager, err := airbot.NewAirBotDataManager(&cloudConn, robotPartID)
	if err != nil {
		logger.Fatal(err)
	}
	//nolint:errcheck
	defer dataManager.Close()

	waypoints, err := waypoint.ReadWaypointsFromFile("./routes/w1-route.csv")
	if err != nil {
		logger.Panic(err)
	}
	logger.Infof("moving to waypoints: %v", waypoints)
	robot, err := client.New(
		context.Background(),
		os.Getenv("ROBOT_ADDRESS"),
		logger,
		client.WithDialOptions(rpc.WithCredentials(rpc.Credentials{
			Type:    utils.CredentialsTypeRobotLocationSecret,
			Payload: os.Getenv("ROBOT_LOCATION_SECRET"),
		})),
	)
	if err != nil {
		logger.Panic(err)
	}
	//nolint:errcheck
	defer robot.Close(ctx)

	logger.Info("successfully connected to the robot")
	a := airbot.NewAirBot(logger, robot, waypoints, dataManager)
	a.Start()
}

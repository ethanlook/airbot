// Package main starts the airbot.
package main

import (
	"context"
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

	waypoints, err := waypoint.ReadWaypointsFromFile("./routes/main-office.csv")
	if err != nil {
		logger.Panic(err)
	}
	logger.Infof("moving to waypoints: %v", waypoints)
	robot, err := client.New(
		context.Background(),
		os.Getenv("ROBOT_LOCATION"),
		logger,
		client.WithDialOptions(rpc.WithCredentials(rpc.Credentials{
			Type:    utils.CredentialsTypeRobotLocationSecret,
			Payload: os.Getenv("ROBOT_SECRET"),
		})),
	)
	if err != nil {
		logger.Panic(err)
	}
	//nolint:errcheck
	defer robot.Close(ctx)
	logger.Info("successfully connected to the robot")
	a := airbot.NewAirBot(logger, robot, waypoints)
	a.Start()
}

// Package main starts the airbot.
package main

import (
	"context"
	"flag"

	"github.com/edaniels/golog"
	"github.com/ethanlook/airbot"
	"github.com/ethanlook/airbot/waypoint"
	"github.com/joho/godotenv"
	"go.viam.com/utils/rpc"

	"go.viam.com/rdk/robot/client"
	"go.viam.com/rdk/utils"
)

var (
	routeName = flag.String("route", "w1", "Route to follow")
)

func main() {
	logger := golog.NewDevelopmentLogger("client")
	ctx := context.Background()

	err := godotenv.Load()
	if err != nil {
		logger.Panic(err)
	}

	waypoints, err := waypoint.ReadWaypointsFromFile("./routes/w1-route.csv")
	if err != nil {
		logger.Panic(err)
	}
	logger.Infof("moving to waypoints: %v", waypoints)
	robot, err := client.New(
		context.Background(),
		"agilexlimo-main.m8n3hqcv6r.viam.cloud",
		logger,
		client.WithDialOptions(rpc.WithCredentials(rpc.Credentials{
			Type:    utils.CredentialsTypeRobotLocationSecret,
			Payload: "k5sa7me1ppx0irjpayrm29pqwb79kf8zanoixr5f62v9ct7u",
		})),
	)
	if err != nil {
		logger.Panic(err)
	}
	//nolint:errcheck
	defer robot.Close(ctx)
	logger.Infof("successfully connected to the robot")
	a := airbot.NewAirBot(logger, robot, waypoints)
	a.Start()
}

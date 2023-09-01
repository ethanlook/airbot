// Package main starts the airbot.
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net"
	"os"

	"github.com/edaniels/golog"
	"github.com/ethanlook/airbot"
	pb "github.com/ethanlook/airbot/proto/v1"
	"github.com/joho/godotenv"
	"go.viam.com/utils/rpc"
	"google.golang.org/grpc"

	"go.viam.com/rdk/robot/client"
	"go.viam.com/rdk/utils"
)

var (
	port         = flag.Int("port", 50051, "The server port")
	errNoRequest = errors.New("missing request")
)

// server is used to implement airbot.AirbotServiceServer.
type server struct {
	pb.UnimplementedAirbotServiceServer

	a      *airbot.AirBot
	logger golog.Logger
}

// Start implements airbot.Start.
func (s *server) Start(ctx context.Context, req *pb.StartRequest) (*pb.StartResponse, error) {
	if req == nil {
		return nil, errNoRequest
	}
	err := s.a.Start(req.GetRoute())
	if err != nil {
		s.logger.Errorf("Error running Start: %w", err)
		return nil, fmt.Errorf("error running Start: %w", err)
	}
	return &pb.StartResponse{}, nil
}

func main() {
	flag.Parse()
	logger := golog.NewDevelopmentLogger("client")
	ctx := context.Background()

	err := godotenv.Load()
	if err != nil {
		logger.Panic(err)
	}

	robotClient, err := client.New(
		context.Background(),
		os.Getenv("ROBOT_LOCATION"),
		logger,
		client.WithDialOptions(rpc.WithCredentials(rpc.Credentials{
			Type:    utils.CredentialsTypeRobotLocationSecret,
			Payload: os.Getenv("ROBOT_SECRET"),
		})),
	)
	if err != nil {
		logger.Fatalf("failed to create robot client")
	}
	defer func() {
		err = robotClient.Close(ctx)
	}()
	logger.Info("successfully connected to the robot")

	a := airbot.NewAirBot(logger, robotClient)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		logger.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterAirbotServiceServer(s, &server{a: a})
	logger.Infof("server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		logger.Fatalf("failed to serve: %v", err)
	}
}

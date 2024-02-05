package delivery

import (
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"net"
	"os"
)

type GrpcServer struct {
	Server *grpc.Server
}

func NewGrpcServer() *GrpcServer {
	return &GrpcServer{
		Server: grpc.NewServer(),
	}
}

func (s *GrpcServer) Start() error {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		return err
	}

	log.Info().Msg("Starting gRPC delivery on :50051")
	return s.Server.Serve(lis)
}

func (s *GrpcServer) WaitForShutdown(signalChan <-chan os.Signal) error {
	sig := <-signalChan
	log.Info().Str("signal", sig.String()).Msg("Shutting down gRPC delivery")

	s.Server.GracefulStop()

	return nil
}

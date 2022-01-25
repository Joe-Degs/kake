package kake

import (
	"log"

	"google.golang.org/grpc"
)

// Server implements the proto.KakeServiceServer interface
type Server struct {
	Config *Config
}

func (s *Server) SelfRegisterService(grpcSrv *grpc.Server) {
	// proto.RegisterKakeServiceServer(grpcSrv, s)
	log.Println("self register called")
}

func DefaultServer() *Server {
	return &Server{
		Config: DefaultConfig(),
	}
}

// NewServer takes a KakeConfig and returns a new Server with the config
func NewServer(cfg *Config) *Server {
	return &Server{
		Config: cfg,
	}
}

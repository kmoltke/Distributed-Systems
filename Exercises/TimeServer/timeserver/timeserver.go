package timeserver

import (
	"golang.org/x/net/context"
	"log"
	"time"
)

type Server struct {
}

func (s *Server) GetTime(ctx context.Context, request *TimeRequest) (*TimeResponse, error) {
	log.Printf("Received request from client: %s", request.ClientName)
	t := time.Now().UnixMicro()
	return &TimeResponse{Time: t}, nil
}

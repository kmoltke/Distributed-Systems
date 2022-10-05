package main

import (
	pb "TimeServer/timeserver"
	"flag"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"time"
)

type Server struct {
	pb.UnimplementedTimeServer        // You need this line if you have a server
	name                       string // Not required but useful if you want to name your server
	port                       string // Not required but useful if your server needs to know what port it's listening to

	//incrementValue int64      // value that clients can increment.
	//mutex          sync.Mutex // used to lock the server to avoid race conditions.
}

// Flags:
var serverName = flag.String("name", "default", "Server name") // set with "-name <name>" in terminal
var port = flag.String("port", "5400", "Server port")          // set with "-port <port>" in terminal

func main() {
	// setLog() // Uncomment this line to log to a log.txt file

	flag.Parse()

	fmt.Println("--- Server is starting ---")

	// start launchServer thread
	go launchServer()

	// Make sure that main is kept alive
	for {
		time.Sleep(time.Second * 5)
	}
}

func launchServer() {
	log.Printf("Server %s: Attempts to create lis on port %s\n", *serverName, *port)

	// Create listener lis tcp on given port or default port 5400
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", *port))
	if err != nil {
		log.Printf("Server %s: Failed to listen on port %s: %v", *serverName, *port, err)
		return
	}

	// Optional options for grpc server:
	var opts []grpc.ServerOption

	// Create pb server (not yet ready to accept requests yet)
	grpcServer := grpc.NewServer(opts...)

	// Make a server instance using the name and port from the flags
	server := &Server{
		name: *serverName,
		port: *port,
	}

	pb.RegisterTimeServer(grpcServer, server)

	log.Printf("Server %s: Listening on port %s\n", *serverName, *port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve %v", err)
	}
}

func (s *Server) GetTime(ctx context.Context, request *pb.TimeRequest) (*pb.TimeResponse, error) {
	log.Printf("Received request from client: %s", request.ClientName)
	t := time.Now().UnixMicro()
	return &pb.TimeResponse{Time: t}, nil
}

// sets the logger to use a log.txt file instead of the console
func setLog() {
	// If a log.txt file exsists, clear the file when a new server is started
	if err := os.Truncate("log.txt", 0); err != nil {
		log.Printf("Failed to truncate: %v", err)
	}

	// Connect to the log file/changes the output of the log informaiton to the log.txt file
	f, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)
}

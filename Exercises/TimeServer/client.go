package main

import (
	pb "TimeServer/timeserver"
	"bufio"
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"os"
	"strings"
	"time"
)

// Flags:
var clientName = flag.String("name", "default-client", "Senders name")
var serverPort = flag.String("port", "5400", "tcp server")

var server pb.TimeClient        // the server
var serverConn *grpc.ClientConn // the server connection

func main() {
	flag.Parse()

	fmt.Println("--- Client App ---")

	// setLog()		// Uncomment this line to log to file

	// Connect to server and close th connection when program closes
	fmt.Println("Attempting to connect to server...")
	connectToServer()
	defer serverConn.Close()

	// Start reading user input
	parseInput()
}

func connectToServer() {
	// Dial Options:
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithBlock(), grpc.WithTransportCredentials(insecure.NewCredentials()))

	// Time out on the connection:
	timeContext, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Dial the server to get a connection:
	log.Printf("Client %s: Attempts to dial on port %s\n", *clientName, *serverPort)
	conn, err := grpc.DialContext(timeContext, fmt.Sprintf(":%s", *serverPort), opts...)
	if err != nil {
		log.Println("Failed to dial: %v", err)
		return
	}

	// TODO: Try to exclude serverConn and only use conn variable
	server = pb.NewTimeClient(conn)
	serverConn = conn
	log.Printf("The connection is: %s\n", conn.GetState().String())
}

func parseInput() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Instruction to user here:")
	fmt.Println("type 'getTime' to request server time")

	for {
		fmt.Print("-> ")

		// Read input into var input and any errors into err
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("Input gave an error: %v", err)
		}
		// Trim input
		input = strings.TrimSpace(input)

		if serverConn.GetState().String() != "READY" {
			//TODO: Try to substitute with log.Fatalf()
			log.Printf("Client %s: Something was wrong with the connection to the server :(", *clientName)
			continue
		}

		//TODO: Event logic goes here:
		if input == "getTime" {
			requestTime()
		}
	}
}

func requestTime() {
	request := &pb.TimeRequest{ClientName: *clientName}

	response, err := server.GetTime(context.Background(), request)
	if err != nil {
		log.Printf("Client %s: No response from server, attempting to reconnect\n", *clientName)
		log.Println(err)
	}

	// Format unix time:
	t := time.UnixMicro(response.Time)

	fmt.Printf("The server time is: %v\n", t)
}

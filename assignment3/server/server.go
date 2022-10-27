package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	proto "simpleGuide/grpc"
	"strconv"

	"google.golang.org/grpc"
)

type Vector_clock struct {
	vector_clock []byte
}

func (local_clock Vector_clock) update(coming_clock Vector_clock, clock_id int) Vector_clock {
	for index, clock := range coming_clock.vector_clock {
		if clock > local_clock.vector_clock[index] {
			local_clock.vector_clock[index] = clock
		}
	}
	local_clock.vector_clock[clock_id]++
	return local_clock
}

var max_vector_size = 10

type Server struct {
	proto.BroadcastServer
	clients            map[int64]Client
	start_client_no    int64
	local_vector_clock Vector_clock
	server_id          int
}

var (
	join      = 0
	left      = 1
	broadcast = 2
)

// Used to get the user-defined port for the server from the command line
var port = flag.Int("port", 0, "server port number")

func main() {
	// Get the port from the command line when the server is run
	flag.Parse()

	//set log file
	logFileName := strconv.Itoa(*port) + ".txt"
	file, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(file)

	log.Println("Log file:" + logFileName + "record started.")
	fmt.Println("I use" + logFileName + "for logging!")

	startServer()

}

type Client struct {
	srv       proto.Broadcast_SendBroadcastServer
	client_id int64
}

func startServer() {

	// Create a new grpc server
	grpcServer := grpc.NewServer()

	// Make the server listen at the given port (convert int port to string)
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(*port))

	if err != nil {
		log.Fatalf("Could not create the server %v", err)
	}
	log.Printf("Started server at port: %d\n", *port)
	fmt.Printf("Started server at port: %d\n", *port)
	server := Server{
		clients:         make(map[int64]Client),
		start_client_no: 1,
		BroadcastServer: proto.UnimplementedBroadcastServer{},
		local_vector_clock: Vector_clock{
			vector_clock: make([]byte, max_vector_size),
		},
		server_id: 0,
	}
	// Register the grpc server and serve its listener
	proto.RegisterBroadcastServer(grpcServer, &server)
	serveError := grpcServer.Serve(listener)
	if serveError != nil {
		log.Fatalf("Could not serve listener")
	}
}

func (s *Server) SendBroadcast(srv proto.Broadcast_SendBroadcastServer) error {
	for {
		msg, err := srv.Recv()
		if err != nil {
			log.Println("receive err:", err)
			return err
		}
		switch msg.Type {
		case int64(join):
			{
				s.clients[msg.ClientId] = Client{srv: srv, client_id: s.start_client_no}
				s.start_client_no++
				msg.Content = fmt.Sprintf("The client %d is joined ", s.clients[msg.ClientId].client_id)
				msg.VectorClock = make([]byte, max_vector_size)
				msg.VectorClock[s.clients[msg.ClientId].client_id] = 1
				log.Println("receive msg from client ", s.clients[msg.ClientId].client_id, "whose clock is ", msg.VectorClock)
				fmt.Println("receive msg from client ", s.clients[msg.ClientId].client_id, "whose clock is ", msg.VectorClock)
				s.Broadcast(msg)
			}
		case int64(left):
			{
				log.Println("receive msg from client ", s.clients[msg.ClientId].client_id, "whose clock is ", msg.VectorClock)
				fmt.Println("receive msg from client ", s.clients[msg.ClientId].client_id, "whose clock is ", msg.VectorClock)
				delete(s.clients, msg.ClientId)
				s.Broadcast(msg)
			}
		case int64(broadcast):
			{
				log.Println("receive msg from client ", s.clients[msg.ClientId].client_id, "whose clock is ", msg.VectorClock)
				fmt.Println("receive msg from client ", s.clients[msg.ClientId].client_id, "whose clock is ", msg.VectorClock)
				s.Broadcast(msg)
			}
		}
	}

}
func (s *Server) Broadcast(msg *proto.Message) {
	coming_vector_clock := Vector_clock{vector_clock: msg.VectorClock}
	s.local_vector_clock.update(coming_vector_clock, s.server_id)
	s.local_vector_clock.update(coming_vector_clock, s.server_id)
	log.Println("Broadcast message: ", msg.Content, " current vector clock is ", s.local_vector_clock.vector_clock)
	fmt.Println("Broadcast message: ", msg.Content, " current vector clock is ", s.local_vector_clock.vector_clock)
	for index, client := range s.clients {
		send_msg := proto.Message{
			Type:        msg.Type,
			Content:     msg.Content,
			VectorClock: s.local_vector_clock.vector_clock,
			ClientId:    client.client_id,
		}
		if err := client.srv.Send(&send_msg); err != nil {
			log.Printf("broadcast client %d err: %v\n", index, err)
		}
	}

}

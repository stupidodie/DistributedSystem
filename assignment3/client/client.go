package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	proto "simpleGuide/grpc"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

type Client struct {
	id                 int64
	local_vector_clock Vector_clock
}

var (
	clientPort = flag.Int("cPort", 0, "client port number")
	serverPort = flag.Int("sPort", 0, "server port number (should match the port used for the server)")
)

func main() {
	// Parse the flags to get the port for the client
	flag.Parse()

	client := Client{
		id:                 -1,
		local_vector_clock: Vector_clock{},
	}
	waitForRequest(&client)
	// Wait for the client (user) to ask for the time

}

func waitForRequest(client *Client) {
	// Connect to the server
	serverConnection, _ := connectToServer()

	// Wait for input in the client terminal
	scanner := bufio.NewScanner(os.Stdin)
	stream, err := serverConnection.SendBroadcast(context.Background())
	stream.Send(&proto.Message{
		Type:        0,
		Content:     "",
		VectorClock: nil,
		ClientId:    int64(*clientPort),
	})
	if err != nil {
		log.Printf("error is %v happen in joining", err)
	}
	go func() {
		for {
			msg, err := stream.Recv()
			if err != nil {
				log.Fatalf("cannot receive %v", err)
			}
			if client.id == -1 {
				client.local_vector_clock.vector_clock = msg.VectorClock
				client.local_vector_clock.update(client.local_vector_clock, int(msg.ClientId))
			} else {
				incoming_vector_clock := Vector_clock{vector_clock: msg.VectorClock}
				client.local_vector_clock.update(incoming_vector_clock, int(msg.ClientId))
			}
			client.id = msg.ClientId
			log.Println("receive msg is ", msg.Content, "current vector clock is", client.local_vector_clock.vector_clock, " id is", msg.ClientId)
		}
	}()
	for scanner.Scan() {

		input := scanner.Text()
		switch input {
		case "exit":
			{
				client.local_vector_clock.update(client.local_vector_clock, int(client.id))
				stream.Send(&proto.Message{Type: 1, Content: fmt.Sprint("The client is going to left, id :", client.id), VectorClock: client.local_vector_clock.vector_clock, ClientId: int64(*clientPort)})
				break
			}
		default:
			{
				client.local_vector_clock.update(client.local_vector_clock, int(client.id))
				stream.Send(&proto.Message{Type: 2, Content: fmt.Sprint("The client id : ", client.id, " want to broadcast ", input), VectorClock: client.local_vector_clock.vector_clock, ClientId: int64(*clientPort)})
			}

		}
	}
}
func connectToServer() (proto.BroadcastClient, error) {
	// Dial the server at the specified port.
	conn, err := grpc.Dial("localhost:"+strconv.Itoa(*serverPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Could not connect to port %d", *serverPort)
	} else {
		log.Printf("Connected to the server at port %d\n", *serverPort)
	}
	return proto.NewBroadcastClient(conn), nil
}

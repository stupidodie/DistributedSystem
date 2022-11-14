package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	p2p "p2p/grpc"
	"strconv"
	"time"

	"google.golang.org/grpc"
)

type peer struct {
	p2p.UnimplementedRingServer
	id      int                    // own id
	clients map[int]p2p.RingClient // map of peer id to client
	ctx     context.Context        // context
}

type Config struct {
	Id    int      `json:"id"`
	Tasks []string `json:"tasks"`
}

const config_file_path = "./config.json"
const start_port = 8000

var configs []Config

const node_numbers = 3

var token_has = false
var Order = 0
var ownPort int

func readConfigFile() {
	byteValue, err := ioutil.ReadFile(config_file_path)
	if err != nil {
		panic(fmt.Sprintln("Cannot read from the config file"))
	}
	configs = make([]Config, 0)
	json.Unmarshal([]byte(byteValue), &configs)
}

func getLocation(ownPort int) int {
	for index, config := range configs {
		if config.Id == ownPort {
			return index
		}
	}
	panic(fmt.Sprintln("cannot find the id for ", ownPort))
}

func getNext(ownPort int) int {
	return configs[((getLocation(ownPort) + 1) % node_numbers)].Id
}

func checkFormer(fromPort int, ownPort int) {
	if configs[((getLocation(ownPort)-1+node_numbers)%node_numbers)].Id != fromPort {
		panic(fmt.Sprintln("The from node is invalid", fromPort))
	}
}

func main() {
	readConfigFile()

	arg1, err := strconv.ParseInt(os.Args[1], 10, 32)
	if err != nil {
		panic(fmt.Sprintln("The parse int is error", err))
	}
	ownPort = int(arg1) + start_port
	ownTasks := configs[getLocation(ownPort)].Tasks
	currentTaskId := 0
	maxTaskNumber := 2
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	peer := &peer{
		id:      ownPort,
		clients: make(map[int]p2p.RingClient),
		ctx:     ctx,
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", ownPort))
	if err != nil {
		log.Fatalf("Failed to listen on port %v: %v", ownPort, err)
	}

	grpcServer := grpc.NewServer()
	p2p.RegisterRingServer(grpcServer, peer)

	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("failed to server %v", err)
		}
	}()

	for i := 0; i < node_numbers; i++ {
		port := configs[i].Id

		if port == int(ownPort) {
			if i == 0 {
				token_has = true
			}
			continue
		}

		var conn *grpc.ClientConn
		conn, err := grpc.Dial(fmt.Sprintf(":%v", port), grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			log.Fatalf("Could not connect: %s", err)
		}

		// Defer means: When this function returns, call this method (meaing, one main is done, close connection)
		defer conn.Close()

		//  Create new Client from generated gRPC code from proto
		peer.clients[port] = p2p.NewRingClient(conn)
	}

	fmt.Println("Press enter to send the token to the nextone - if you has")
	go func() {
		for {
			if token_has {
				fmt.Println("Know I have the token, it will send to another after 3 seconds")
				time.Sleep(3 * time.Second)
				token_has = false
				peer.sendToNextOne()
			}
		}
	}()
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		if token_has {
			if currentTaskId <= maxTaskNumber {
				fmt.Println("The task is ", ownTasks[currentTaskId], "Order is ", Order)
			} else {
				fmt.Println("current is no task for this node, just handout the token")
			}
			currentTaskId += 1
			// token_has=false
			// send to the next one
			// peer.sendToNextOne()
		} else {
			fmt.Println("current node does not hold token, so just waiting")
		}

	}

}

// The peer receives a ping from another peer:
func (p *peer) HandNext(ctx context.Context, req *p2p.MSG) (*p2p.Reply, error) {
	id := req.Id
	order := req.Order
	Order = int(order) + 1
	checkFormer(int(id), ownPort)
	token_has = true
	return &p2p.Reply{Msg: fmt.Sprintln("Receive token", order, " from ", id)}, nil
}

func (p *peer) sendToNextOne() {
	next_node := getNext(ownPort)
	request := &p2p.MSG{Id: int32(ownPort), Order: int32(Order)}
	reply, err := p.clients[next_node].HandNext(p.ctx, request)
	if err != nil {
		log.Fatalf("Error when sending ping to peer %v: %v", next_node, err)
	}
	fmt.Println("Send msg to the next node port is ", next_node, "and receive the reply", reply.Msg)
}

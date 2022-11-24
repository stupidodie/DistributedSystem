package main

import (
	gp "assignment5/grpc"
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)
const config_file_path = "./config.json"

type Config struct {
	Port int `json:"port"`
}

var configs []Config

type Info struct {
	servers []gp.TradeClient
	ctx     context.Context // context
}

func main() {
	readConfigFile()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	info := &Info{servers: make([]gp.TradeClient, len(configs)), ctx: ctx}

	for i, config := range configs {
		var conn *grpc.ClientConn
		conn, err := grpc.Dial(fmt.Sprintf(":%v", config.Port), grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			log.Fatalf("Could not connect: %s", err)
		}
		// Defer means: When this function returns, call this method (meaing, one main is done, close connection)
		defer conn.Close()
		//  Create new Client from generated gRPC code from proto
		info.servers[i] = gp.NewTradeClient(conn)
	}

	scanner := bufio.NewScanner(os.Stdin)
outer:
	for scanner.Scan() {
		input := scanner.Text()
		switch input {
		case "q":
			break outer

		case "quit":
			break outer

		default:
			price,err := strconv.ParseInt(input, 10, 32)
			if err != nil {
				log.Fatalf("parse error %v", err)
			}
			var ack int32
			var higest_price int32
			var is_finished bool
			for _, v := range info.servers {
				request := &gp.Price{Price: int32(price)}
				reply, err := v.Bid(info.ctx, request)
				if err != nil {
					continue
				}
				ack=reply.Ack
				emp := &emptypb.Empty{}
				outcome, err1 := v.Result(info.ctx, emp)
				if err1 != nil {
					continue
				}
				higest_price=outcome.Price
				is_finished=outcome.IsFinished
			}
			log.Println("ack is ",ack)
			log.Println("reply price is ", higest_price, " is finished is ",is_finished)
		}
		
	}

}
func readConfigFile() {
	byteValue, err := ioutil.ReadFile(config_file_path)
	if err != nil {
		panic(fmt.Sprintln("Cannot read from the config file"))
	}
	configs = make([]Config, 0)
	json.Unmarshal([]byte(byteValue), &configs)
}

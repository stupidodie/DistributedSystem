package main

import (
	gp "assignment5/grpc"
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)
var max_bid_number=5
type Server struct{
	gp.UnimplementedTradeServer
	highest_price int 
	is_finished bool
}
const success=0
const fail=-1
const exception=1
// Used to get the user-defined port for the server from the command line
var port = flag.Int("port", 0, "server port number")
func main(){
	startServer()
	fmt.Println("das")
	
}
func startServer() {

	// Create a new grpc server
	grpcServer := grpc.NewServer()

	agr1, err := strconv.ParseInt(os.Args[1], 10, 32)
	port:=int(agr1)
	if err!=nil{
		panic(fmt.Sprintln("The parse int is error",err))
	}
	// Make the server listen at the given port (convert int port to string)
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d",port))

	if err != nil {
		log.Fatalf("Could not create the server %v", err)
	}
	log.Printf("Started server at port: %d\n", port)
	
	server:=Server{highest_price: -1, is_finished: false}
	
	// Register the grpc server and serve its listener
	gp.RegisterTradeServer(grpcServer,&server)
	serveError := grpcServer.Serve(listener)
	if serveError != nil {
		log.Fatalf("Could not serve listener")
	}
}

func (s *Server) Bid(ctx context.Context,req *gp.Price) (*gp.Ack,error){
	max_bid_number-=1
	if max_bid_number<0{
		s.is_finished=true
	}
	if s.is_finished{
		log.Println("Receive bid but current is closed")
		return &gp.Ack{Ack:exception},nil
	}else{
		if req.Price > int32(s.highest_price){
			log.Println("Bid successfully price, current new higest price is ",req.Price)
			s.highest_price=int(req.Price)
			return &gp.Ack{Ack: success},nil
		}else{
			log.Println("Bid failed price is ",req.Price)
			return &gp.Ack{Ack: fail},nil
		}
	}	
}


func (s* Server) Result(ctx context.Context,emp *emptypb.Empty)(*gp.Outcome,error){
	return &gp.Outcome{Price:int32(s.highest_price),IsFinished:s.is_finished },nil
}
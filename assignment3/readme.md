# Time Request gRPC Example Guide

-[How to Run](#how-to-run)

## How to Run
For example, if you want to simulate three clients communicating with each other, try:
1. Run the server: `go run server/server.go -port 5454`.
2. In a different terminal, run the client: `go run client/client.go -cPort 8083 -sPort 5454`.
3. In a different terminal, run the client: `go run client/client.go -cPort 8084 -sPort 5454`.
4. In a different terminal, run the client: `go run client/client.go -cPort 8085 -sPort 5454`.
5. In the client terminal input something and press enter. You should now get the message from each other.
6. In order to let a client leave the channel, simply type "exit" and other client will know you are offline.

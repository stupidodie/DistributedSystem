package main

import (
	common "common"
	"encoding/json"
	"fmt"
	"net"
)

const init_seq = 5000

func handleConnection(conn net.Conn) {
	defer conn.Close()
	var msg string
	buf := [512]byte{}
	packet := common.NewTCP_Packet()
	for {
		n, err := conn.Read(buf[:])
		if err != nil {
			fmt.Println("recv failed, err:", err)
			return
		}
		err = json.Unmarshal(buf[:n], &packet)
		if err != nil {
			fmt.Println("Unmarshal failed, err:", err, 24)
			return
		}
		switch packet.OPCODE {
		case common.SYN:
			fmt.Println("Start receive SYN,ACK:",packet.ACK,"SEQ:",packet.SEQ)
			packet.OPCODE = common.SYN_ACK
			packet.ACK = packet.SEQ + 1
			packet.SEQ = init_seq
			buf, err := json.Marshal(packet)
			if err != nil {
				fmt.Println("Failed to Marshal", err)
				return
			}
			fmt.Println("Start Send ACK:",packet.ACK,"SEQ:",packet.SEQ)
			conn.Write(buf)
		case common.ACK:
			fmt.Println("Receive the Ack",packet.ACK,"SEQ:",packet.SEQ)
		case common.ESTABLISHED:
			if packet.IS_EXIT {
				fmt.Println("The msg is:", msg)
				return
			}
			msg+= (string)(packet.DATA[:])
			packet.SEQ++
			buf, err := json.Marshal(packet)
			if err != nil {
				fmt.Println("Failed to Marshal", err)
				return
			}
			conn.Write(buf)

		}

	}

}

func main() {
	ls, err := net.Listen(common.PROTOCOL_TYPE, common.ADDRESS)
	if err != nil {
		panic(fmt.Sprintln("Listen failed error: ", err))
	}
	for {
		conn, err := ls.Accept()
		if err != nil {
			fmt.Println("accept failed, err:", err)
			continue
		}
		go handleConnection(conn)
	}
}

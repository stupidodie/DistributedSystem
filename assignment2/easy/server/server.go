package main

import (
	common "common"
	"encoding/json"
	"fmt"
	"net"
	"time"
)

const init_seq = 5000

func handleConnection(conn net.Conn) {
	initSeq := init_seq
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
			fmt.Println("Received a SYN,ACK for first hand-shake:", packet.ACK, "SEQ:", packet.SEQ)
			newPacket := common.NewTCP_Packet()
			newPacket.OPCODE = common.SYN_ACK
			newPacket.ACK = packet.SEQ + 1
			newPacket.SEQ = initSeq
			newPacket.DATA = [common.MAX_DATA_SIZE]byte{0}
			initSeq++

			buf, err := json.Marshal(newPacket)
			if err != nil {
				fmt.Println("Failed to Marshal", err)
				return
			}
			fmt.Println("Start Send SYN_ACK for second hand-shake, ACK is:", newPacket.ACK, "SEQ:", newPacket.SEQ)
			conn.Write(buf)
		case common.ACK:
			fmt.Println("Receive the Ack for the third hand-shake, ACK is:", packet.ACK, "SEQ:", packet.SEQ)
		case common.ESTABLISHED:
			if packet.IS_EXIT {
				fmt.Println("The msg is:", msg)
				return
			}
			msg += (string)(packet.DATA[:])
			fmt.Println("Received data with seq:", packet.SEQ)
			// packet.SEQ++
			newPacket := common.NewTCP_Packet()
			newPacket.SEQ = initSeq
			initSeq++
			newPacket.OPCODE = common.ESTABLISHED
			newPacket.DATA = packet.DATA
			newPacket.ACK = packet.SEQ + 1
			newPacket.DATA = [common.MAX_DATA_SIZE]byte{0}
			//fmt.Println("Sending:", newPacket)
			fmt.Println("Sending the Ack with packet number:", newPacket.SEQ, "ACK:", newPacket.ACK)
			buf, err := json.Marshal(newPacket)
			if err != nil {
				fmt.Println("Failed to Marshal", err)
				return
			}
			time.Sleep(1 * time.Second)
			conn.Write(buf)
			time.Sleep(1 * time.Second)

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

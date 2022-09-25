package main

import (
	"bytes"
	common "common"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"
)

const init_seq = 5000

func handleConnection(conn net.Conn) {
	log.SetFlags(log.Lshortfile)
	initSeq := init_seq
	defer conn.Close()
	var msg string
	buf := [512]byte{}
	packet := common.NewTCP_Init_Packet()
	for {
		n, err := conn.Read(buf[:])
		if err != nil {
			log.Println("recv failed, err:", err)
			return
		}
		err = json.Unmarshal(buf[:n], &packet)
		if err != nil {
			log.Println("Unmarshal failed, err:", err)
			return
		}
		switch packet.OPCODE {
		case common.SYN:
			fmt.Println("Received a SYN,","SEQ:", packet.SEQ,"ACK for first hand-shake:", packet.ACK )
			newPacket:=common.NewTCP_Specific_Packet(initSeq,packet.SEQ + 1,common.SYN_ACK)
			initSeq++
			buf, err := json.Marshal(newPacket)
			if err != nil {
				fmt.Println("Failed to Marshal", err)
				return
			}
			fmt.Println("Start Send SYN_ACK for second hand-shake,","SEQ:", newPacket.SEQ,"ACK is:", newPacket.ACK, )
			conn.Write(buf)
		case common.ACK:
			fmt.Println("Receive the Ack for the third hand-shake,","SEQ:", packet.SEQ,"ACK is:", packet.ACK, )
		case common.ESTABLISHED:
			if packet.IS_EXIT {
				fmt.Println("The message is:", msg)
				return
			}
			//We filter the 0 byte to get the original message
			msg += (string)((bytes.Split(packet.DATA[:],[]byte{0}))[0])
			fmt.Println("Received data with seq:", packet.SEQ)
			newPacket:=common.NewTCP_Specific_Packet(initSeq,packet.SEQ + 1,common.ESTABLISHED)
			buf, err := json.Marshal(newPacket)
			if err != nil {
				log.Println("Failed to Marshal", err)
				return
			}
			time.Sleep(2 * time.Second)
			fmt.Println("Sending the Ack with packet number:", newPacket.SEQ, "ACK:", newPacket.ACK)
			conn.Write(buf)
			time.Sleep(2 * time.Second)
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

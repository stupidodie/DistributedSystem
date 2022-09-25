package main

import (
	common "common"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"
)

// const msg = "TEST MSG!"
const msg = "Here is the message that we want to send"

func main() {
	log.SetFlags(log.Lshortfile)
	conn, err := net.Dial(common.PROTOCOL_TYPE, common.ADDRESS)
	if err != nil {
		fmt.Println("err :", err)
		return
	}
	defer conn.Close()
	packet := common.NewTCP_Init_Packet() //create a new packet to establish connection
	init_seq := packet.SEQ           //we set default sequence number to be 1000
	is_sent := false                 //this state represent if our client has sent all the packets or not
	var last_seq int

	msg_byte_array := []byte(msg)
	msg_byte_array_size := len(msg_byte_array)
	var index int = 0 // for packet index
	buf := [512]byte{} //create a buffer to send and recieve the message
	for {
		switch packet.OPCODE {
		case common.UNSTART:
			packet.OPCODE = common.SYN
			buf, err := json.Marshal(packet)
			if err != nil {
				log.Panicln("Failed to Marshal", err)
			}
			fmt.Println("Start to send the SYN for first hand-shake:", packet.SEQ)
			conn.Write(buf) //staring first hand-shake
		case common.SYN:
			n, err := conn.Read(buf[:])
			if err != nil {
				log.Panicln("recv failed, err:", err)
			}
			err = json.Unmarshal(buf[:n], &packet)
			if err != nil {
				log.Panicln("Unmarshal failed, err:", err)
			}
			if packet.OPCODE == common.SYN_ACK && packet.ACK == (init_seq+1) { //make sure the sequence number is matched
				fmt.Println("Receive the SYN_ACK for second hand-shake", packet.ACK, "SEQ:", packet.SEQ)
				packet.OPCODE = common.ACK
				packet.ACK = packet.SEQ + 1
				packet.SEQ = init_seq + 1
				buf, err := json.Marshal(packet)
				if err != nil {
					log.Panicln("Failed to Marshal", err)
				}
				fmt.Println("Send the ACK for the third hand shake, ACK:", packet.ACK, "seq:", packet.SEQ)
				conn.Write(buf)
				time.Sleep(1 * time.Second) //to avoid sticking packets and create a ideal condition for our experiment
				packet.OPCODE = common.ESTABLISHED //connection build
			}
		case common.ESTABLISHED:
			if !is_sent { //keep sending data packets until the end of the message
				if index+common.MAX_DATA_SIZE < msg_byte_array_size { //if it is not the last packet
					copy(packet.DATA[:], msg_byte_array[index:(index+common.MAX_DATA_SIZE)])
					index +=common.MAX_DATA_SIZE
				} else { //if it is the last packet
					copy(packet.DATA[0:(msg_byte_array_size-index)], msg_byte_array[index:])
					tmp := [common.MAX_DATA_SIZE]byte{0}
					copy(packet.DATA[(msg_byte_array_size-index):], tmp[:])
					is_sent = true //after finish sending all the packets, set state to true
				}
				packet.SEQ++

				buf, err := json.Marshal(packet)
				if err != nil {
					log.Panicln("Failed to Marshal", err)
				}
				last_seq = packet.SEQ
				time.Sleep(1 * time.Second)
				fmt.Println("Sending data with seq:", packet.SEQ)
				conn.Write(buf[:])
				time.Sleep(1 * time.Second)
				break
			}
			fmt.Println("repeat listening")
			n, err := conn.Read(buf[:])
			if err != nil {
				log.Panicln("recv failed, err:", err)
			}
			err = json.Unmarshal(buf[:n], &packet)
			if err != nil {
				log.Panicln("Unmarshal failed, err:", err)
			}
			fmt.Println("Receive the ACK for packet",packet.ACK," seq:", packet.SEQ)

			if packet.ACK == last_seq+1 { // already get the last ack, job finished, then send exit packet to server
				fmt.Println("Receive the last ack, seq:", packet.SEQ)
				packet.SEQ++
				packet.IS_EXIT = true
				buf, err := json.Marshal(packet)
				if err != nil {
					log.Panicln("Failed to Marshal", err)
				}
				time.Sleep(1 * time.Second)
				conn.Write(buf)
				return
			}
		}
	}
}

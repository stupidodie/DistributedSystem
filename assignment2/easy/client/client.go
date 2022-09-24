package main

import (
	common "common"
	"encoding/json"
	"fmt"
	"net"
	"time"
)

const msg = "Here is the message that we want to send"

func main() {
	conn, err := net.Dial(common.PROTOCOL_TYPE, common.ADDRESS)
	if err != nil {
		fmt.Println("err :", err)
		return
	}
	defer conn.Close()
	packet := common.NewTCP_Packet()
	init_seq := packet.SEQ
	is_sent := false
	var last_seq int
	
	for {
		buf := [512]byte{}
		switch packet.OPCODE {
		case common.UNSTART:
			packet.OPCODE = common.SYN
			buf, err := json.Marshal(packet)
			if err != nil {
				fmt.Println("Failed to Marshal", err)
				return
			}
			fmt.Println("Start to send the information seq:",packet.SEQ)
			conn.Write(buf)
		case common.SYN:
			n, err := conn.Read(buf[:])
			if err != nil {
				fmt.Println("recv failed, err:", err,40)
				return
			}
			err = json.Unmarshal(buf[:n], &packet)
			if err != nil {
				fmt.Println("Unmarshal failed, err:", err,58)
				return
			}
			if packet.OPCODE == common.SYN_ACK && packet.ACK == (init_seq+1) {
				fmt.Println("Receive the ACK",packet.ACK,"SEQ:",packet.SEQ)
				packet.OPCODE = common.ACK
				packet.ACK = packet.SEQ + 1
				packet.SEQ = init_seq + 1
				buf, err := json.Marshal(packet)
				if err != nil {
					fmt.Println("Failed to Marshal", err)
					return
				}
				time.Sleep(1000)
				fmt.Println("Send the ACK:",packet.ACK,"seq:",packet.SEQ)
				conn.Write(buf)
				time.Sleep(1000)
				
				packet.OPCODE = common.ESTABLISHED
			}
		case common.ESTABLISHED:
			if !is_sent {
				msg_byte_array:=[]byte(msg)
				msg_byte_array_size:=len(msg_byte_array)	
				for i := 0; i < msg_byte_array_size; i += common.MAX_DATA_SIZE {
					if i+common.MAX_DATA_SIZE< msg_byte_array_size{
						copy(packet.DATA[:], msg_byte_array[i:(i+common.MAX_DATA_SIZE)])
					}else{
						copy(packet.DATA[0:(msg_byte_array_size-i)], msg_byte_array[i:])	
						tmp:=[common.MAX_DATA_SIZE]byte{0}
						copy(packet.DATA[(msg_byte_array_size-i):],tmp[:])
					}
					packet.SEQ++
					buf, err := json.Marshal(packet)
					if err != nil {
						fmt.Println("Failed to Marshal", err)
						return
					}
					last_seq = packet.SEQ
					time.Sleep(1000)
					conn.Write(buf[:])
					time.Sleep(1000)	
				}
				is_sent = true
				break
			}
			n, err := conn.Read(buf[:])
			if err != nil {
				fmt.Println("recv failed, err:", err,78)
				return
			}
			err = json.Unmarshal(buf[:n], &packet)
			if err != nil {
				fmt.Println("Unmarshal failed, err:", err,97)
				return
			}
			if packet.SEQ == last_seq+1 {
				packet.SEQ++
				packet.IS_EXIT = true
				buf, err := json.Marshal(packet)
				if err != nil {
					fmt.Println("Failed to Marshal", err)
					return
				}
				time.Sleep(1000)	
				conn.Write(buf)
				time.Sleep(1000)	
				return
			}
		}
	}
}

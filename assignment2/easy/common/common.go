package common

const (
	UNSTART = iota
	SYN
	SYN_ACK
	ACK
	ESTABLISHED
)
const MAX_DATA_SIZE = 64

type TCP_Packet struct {
	OPCODE  int `json:"opcode"`
	SEQ     int `json:"seq"`
	IS_EXIT bool  `json:"is_exit"`
	DATA    [MAX_DATA_SIZE]byte `json:"data"`
	ACK     int `json:"ack"`
}


func NewTCP_Packet() TCP_Packet{
	return TCP_Packet{
		OPCODE: UNSTART,
		SEQ: 1000,
		IS_EXIT: false,
		DATA: [MAX_DATA_SIZE]byte{0},
		ACK: 0,
	}
}

const ADDRESS = "127.0.0.1:20000"
const PROTOCOL_TYPE = "tcp"
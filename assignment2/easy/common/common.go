package common

const (
	UNSTART = iota
	SYN
	SYN_ACK
	ACK
	ESTABLISHED
)
const MAX_DATA_SIZE = 16

type TCP_Packet struct {
	OPCODE  int                 `json:"opcode"`
	SEQ     int                 `json:"seq"`
	IS_EXIT bool                `json:"is_exit"`
	DATA    [MAX_DATA_SIZE]byte `json:"data"`
	ACK     int                 `json:"ack"`
}

func NewTCP_Init_Packet() TCP_Packet {
	return TCP_Packet{
		OPCODE:  UNSTART,
		SEQ:     1000,
		IS_EXIT: false,
		DATA:    [MAX_DATA_SIZE]byte{0},
		ACK:     0,
	}
}
func NewTCP_Specific_Packet(Seq int,Ack int,opcode int) TCP_Packet {
	return TCP_Packet{
		OPCODE: opcode,
		SEQ:     Seq,
		IS_EXIT: false,
		DATA:    [MAX_DATA_SIZE]byte{0},
		ACK:     Ack,
	}
}




const ADDRESS = "127.0.0.1:20000"
const PROTOCOL_TYPE = "tcp"

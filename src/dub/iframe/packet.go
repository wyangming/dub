package iframe

import "encoding/binary"

const (
	InvliadPacketMainID = uint16(0xFFFF)
)
const (
	PacketCommandSize = 4 //MainId(uint16) + SubId(uint16)
)

//普通封包
type Packet struct {
	MainId uint16 //主命令ID
	SubId  uint16 //子命令ID
	Data   []byte //数据
}

func (self *Packet) ToByteArray() []byte {

	buffer := make([]byte, PacketCommandSize)

	binary.BigEndian.PutUint16(buffer[0:2], self.MainId)

	binary.BigEndian.PutUint16(buffer[2:4], self.SubId)

	buffer = append(buffer[0:PacketCommandSize], self.Data...)

	return buffer
}

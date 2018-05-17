package frame

import (
	"encoding/binary"
	"errors"
)

const (
	//代表错误的命令
	InvliadPacketMainID = uint16(0xFFFF)
	//代表网络命令
	NetPacketMainID = uint16(0x0000)
	//需要与NetPacketMainID一块使用代表的是心跳连接
	NetPacketHeatBateSubID = uint16(0x0000)
)

type Packet struct {
	MainId uint16 //主命令
	SubId  uint16 //子命令
	Data   []byte //数据
}

//转为字节流
func (this *Packet) ToByteArray() []byte {
	dateCount := len(this.Data)

	buffer := make([]byte, 0, 4+dateCount)

	//uint16是16位的数字一个占两个字节
	binary.BigEndian.PutUint16(buffer[0:2], this.MainId)
	binary.BigEndian.PutUint16(buffer[2:4], this.SubId)

	buffer = append(buffer[0:4], this.Data...)
	return buffer
}

//把字节流转为想要的包
func ToDecryptPacket(data []byte) (p *Packet, err error) {
	//如果长度小于4则说明发送的消息格式不对
	if len(data) < 4 {
		return nil, errors.New("Packet length is too short")
	}
	p = &Packet{}
	p.MainId = uint16(binary.BigEndian.Uint16(data[0:2]))
	p.SubId = uint16(binary.BigEndian.Uint16(data[2:4]))
	p.Data = data[4:]

	return p, nil
}

//创建一个包信息
func MakePacket(mainId, subId uint16, data []byte) *Packet {
	return &Packet{
		MainId: mainId,
		SubId:  subId,
		Data:   data,
	}
}

package frame

import (
	"encoding/binary"
	"io"
	"net"
)

const (
	//包大小的位数
	PACKETSTREAMHEAD = uint64(8)
)

//封包流
type PacketStream interface {
	Read() (*Packet, error)
	Write(pkt *Packet) error
	Close() error
	Raw() net.Conn
}

type pntStream struct {
	conn net.Conn
}

func (this *pntStream) Read() (p *Packet, err error) {
	headdata := make([]byte, PACKETSTREAMHEAD)
	if _, err = io.ReadFull(this.conn, headdata); err != nil {
		return nil, err
	}
	data := make([]byte, uint64(binary.BigEndian.Uint64(headdata))-PACKETSTREAMHEAD)
	if _, err = io.ReadFull(this.conn, data); err != nil {
		return nil, err
	}
	p, err = ToDecryptPacket(data)
	if err != nil {
		return nil, err
	}

	return p, nil
}
func (this *pntStream) Write(pkt *Packet) (err error) {
	pktData := pkt.ToByteArray()
	dataLen := uint64(len(pktData)) + PACKETSTREAMHEAD
	data := make([]byte, 0, dataLen)
	binary.BigEndian.PutUint64(data[0:PACKETSTREAMHEAD], dataLen)
	data = append(data[:PACKETSTREAMHEAD], pktData...)

	if _, err = this.conn.Write(data); err != nil {
		return err
	}
	return nil
}
func (this *pntStream) Close() error {
	return this.conn.Close()
}

//返回最原始的连接
func (this *pntStream) Raw() net.Conn {
	return this.conn
}

//创建封包流
func NewPacketStream(conn net.Conn) PacketStream {
	return &pntStream{
		conn: conn,
	}
}

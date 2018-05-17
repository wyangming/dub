package iframe

import (
	"net"
	"sync"
	"time"
)

const (
	NetPacketSecond = 8 //发送心中包的间隔
)

type SocketConnector struct {
	conn             net.Conn                          //连接
	autoReconnectSec int                               //重连间隔，0为不重连
	cloneSignal      chan bool                         //关闭信号
	stream           PacketStream                      //读写流
	writeChan        interface{}                       //写入通道
	iskeep           bool                              //是否保持长连
	ip               string                            //连接服务器的ip
	endSync          sync.WaitGroup                    //同步锁
	isdone           bool                              //标志是否断开
	CallBack         func(uint16, uint16, []byte) bool //接收到包后的处理函数
	ticker           *time.Ticker                      //发送心中包的定时器
}

//启动连接服务器
func (s *SocketConnector) Start(ip string) error {
	s.ip = ip
	return s.connect()
}

//断开连接
func (s *SocketConnector) Stop() bool {
	//关闭写入通道
	s.writeChan <- false
	//通知线程接收
	s.endSync.Done()
	return true
}

//发送命令
func (s *SocketConnector) Send(mainId uint16, subId uint16, data []byte) error {
	pkt := MakePacket(mainId, subId, data)
	if s.isdone {
		return errors.New("connector has done!")
	} else {
		s.writeChan <- pkt
	}

	return nil
}

func (s *SocketConnector) sendThread() {
	for {
		switch pkt := (<-s.writeChan).(type) {
		// 封包
		case *Packet:
			if err := s.stream.Write(pkt); err != nil {
				goto exist_loop
			}
		case bool:
			goto exist_loop
		}
	}

exist_loop:
	// 关闭socket,触发读错误, 结束读循环
	s.stream.Close()
	// 通知接收线程ok
	s.endSync.Done()
}

//接收服务器返回的信息
func (s *SocketConnector) recvThread() {
	for {
		// 从Socket读取封包
		pk, err := s.stream.Read()

		if err != nil {
			break
		}

		//心跳包
		if pk.MainId == NetPacketMainID && pk.SubId == NetPacketHeatBateSubID {
			continue
		}

		ret := s.CallBack(pk.MainId, pk.SubId, pk.Data)
		if false == ret {
			break
		}
	}

	if !s.isdone {
		// 通知关闭写协程
		s.writeChan <- true
		// 通知接收线程ok
		s.endSync.Done()
	}
}

func (s *SocketConnector) existThread() {
	// 布置接收和发送2个任务
	s.endSync.Add(2)
	// 等待2个任务结束
	s.endSync.Wait()

	if s.iskeep {
		for {
			err := s.connect()
			if err == nil {
				break
			}
		}

	}
}

//连接服务器
func (s *SocketConnector) connect() error {
	tcpAddr, err := net.ResolveTCPAddr("tcp", s.ip)
	if err != nil {
		return err
	}

	s.conn, err = net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return err
	}

	s.stream = NewPacketStream(s.conn)

	// 退出线程
	go s.existThread()

	// 接收线程
	go s.recvThread()

	// 发送线程
	go s.sendThread()

	//发送心跳包
	go s.netPacket()

	s.isdone = false

	return nil
}

//发送心跳包
func (s *SocketConnector) netPacket() {
	s.ticker = time.NewTicker(NetPacketSecond * time.Second)
	for _ = range s.ticker.C {
		s.Send(NetPacketMainID, NetPacketHeatBateSubID, nil)
	}
}

//新建连接
func NewConnector() *SocketConnector {
	self := &SocketConnector{
		writeChan:   make(chan interface{}),
		cloneSignal: make(chan bool),
		isdone:      true,
		iskeep:      true,
	}

	return self
}

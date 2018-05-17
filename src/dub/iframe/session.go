package frame

import (
	"dub/common"
	"errors"
	"net"
	"sync"
	"time"
)

const (
	SessionDeadSecond = 10 //session的超时时间，以秒为单位
)

type Session struct {
	sockid        uint64
	conn          net.Conn
	OnClose       func(common.ISession)
	CallBack      func(uint16, uint16, []byte, common.ISession) //接收到包后的处理函数
	stream        PacketStream
	endSync       sync.WaitGroup
	needStopWrite bool                   //是否需要主动断开协程
	isDone        bool                   //是否已经关闭
	sendList      *PacketList            //将发包改为发送列表
	data          map[string]interface{} //session中存的数据
}

//发包
func (this *Session) Send(mainId uint16, subId uint16, data []byte) error {
	pkt := MakePacket(mainId, subId, data)
	if this.isDone {
		return errors.New("session has done")
	}
	this.sendList.Add(pkt)

	return nil
}

//断开
func (this *Session) Close(closeType common.CloseType) {
	if !this.isDone {
		this.isDone = true
		// 通知关闭写协程
		this.sendList.Close()
	}
}

//标识ID
func (this *Session) ID() uint64 {
	return this.sockid
}

func (this *Session) GetRemoteAddr() string {
	ipstring, _, _ := net.SplitHostPort(this.conn.RemoteAddr().String())
	return ipstring
}
func (this *Session) SetSessionData(key string, val interface{}) {
	if val != nil&len(key) > 0 {
		this.data[key] = val
	}
}
func (this *Session) GetSessionData(key string) interface{} {
	if len(key) > 0 {
		val, ok := this.data[key]
		if !ok {
			return val
		}
	}
	return nil
}
func (this *Session) DelSessionData(key string) {
	if len(key) > 0 {
		_, ok := this.data[key]
		if ok {
			delete(this.data, key)
		}
	}
}

func (this *Session) existThread() {
	// 等待2个任务结束
	this.endSync.Wait()

	if this.OnClose != nil {
		this.sendList.Close()
		this.OnClose(this)
	}
}

func (this *Session) recvThread() {
	for {
		pk, err := this.stream.Read()
		if err != nil {
			break
		}

		//判断是否为心跳检测包
		if pk.MainId == NetPacketMainID && pk.SubId == NetPacketHeatBateSubID {
			self.conn.SetDeadline(time.Now().Add(time.Duration(SessionDeadSecond) * time.Second))
			this.Send(NetPacketMainID, NetPacketHeatBateSubID, nil)
		} else {
			if this.CallBack != nil {
				this.CallBack(pk.MainId, pk.SubId, pk.Data, this)
			}
		}
	}

	this.isDone = true
	if this.needStopWrite {
		this.stream.Close()
	}

	// 通知接收线程ok
	this.endSync.Done()
}

func (this *Session) sendThread() {
	var sendList []*Packet
	for {
		sendList = sendList[0:0]

		//复制出队列
		packetList := this.sendList.BeginPick()
		sendList = append(sendList, packetList...)
		this.sendList.EndPick()

		willExit := false
		//写队列的数据
		for _, p := range sendList {
			if p.MainId == InvliadPacketMainID {
				willExit = true
			} else if err := this.stream.Write(p); err != nil {
				willExit = true
				break
			}
		}

		if willExit {
			goto exitSendLoop
		}
	}

exitSendLoop:
	this.isDone = true
	this.needStopWrite = false
	this.stream.Close()

	this.endSync.Done()
}

func NewSession(sockid uint64, tcpConn net.Conn) *Session {
	self := &Session{
		sockid:        sockid,
		conn:          tcpConn,
		needStopWrite: true,
		isDone:        false,
		sendList:      NewPacketList(),
	}
	//设置超时时间
	self.conn.SetDeadline(time.Now().Add(time.Duration(SessionDeadSecond) * time.Second))
	self.stream = NewPacketStream(self.conn)
	self.endSync.Add(2)

	// 退出线程
	go self.existThread()
	// 接收线程
	go self.recvThread()
	// 发送线程
	go self.sendThread()

	return self
}

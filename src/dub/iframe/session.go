package iframe

import (
	"net"
	"dub/common"
)

type Session struct {
	socketId uint64
	conn net.Conn
	writerChan chan interface{}
	OnClose func(session common.ISession)

}

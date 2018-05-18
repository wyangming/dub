package frame

import (
	"net"
	"sync"
	"time"
	"dub/utils"
)

type ListenerTcp struct {
	listener  net.Listener //监听
	running   bool         //运行状态
	sessions  map[uint64]*Session
	fdcounter uint64                            //统计共有多少个session连接
	CallBack  func(uint16, uint16, []byte) bool //接收到包后的处理函数
	mutex     sync.Mutex                        //同步锁
	ticker    time.Ticker                       //定时检测任务
	log       *utils.Logger                     //日志信息
}

func NewListener(callBack func(uint16, uint16, []byte) bool) *ListenerTcp {
	self := &ListenerTcp{
		fdcounter: 0,
		log:       utils.NewLogger(),
		CallBack:  callBack,
	}
	return self
}

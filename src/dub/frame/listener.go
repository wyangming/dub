package frame

import (
	"dub/common"
	"dub/utils"
	"net"
	"sync"
	"time"
)

const (
	ListenerClearSessionSec = 30 //服务器清理过期session的时间
)

type ListenerTcp struct {
	addr       string       //服务监听的地址与端口
	listener   net.Listener //监听
	running    bool         //运行状态
	sessions   map[uint64]*Session
	fdcounter  uint64                                             //统计共有多少个session连接
	CallBack   func(uint16, uint16, []byte, common.ISession) bool //接收到包后的处理函数
	OnShutdown func(common.ISession) bool                         //断开命令
	mutex      sync.Mutex                                         //同步锁
	ticker     *time.Ticker                                       //定时检测任务
	log        *utils.Logger                                      //日志信息
}

//停止服务
func (l *ListenerTcp) Stop() {
	if l.running {
		l.running = false
		l.ticker.Stop()
		l.listener.Close()
	}
}

//清理过期的session
func (l *ListenerTcp) clearSession() {
	for _ = range l.ticker.C {
		var deadSesVec []*Session

		l.mutex.Lock()

		for _, ses := range l.sessions {
			if ses != nil && ses.ConnState() {
				deadSesVec = append(deadSesVec, ses)
			}
		}

		l.mutex.Unlock()

		for _, ses := range deadSesVec {
			if l.OnShutdown!=nil {
				l.OnShutdown(ses)
			}
			if ses.OnClose != nil {
				ses.OnClose(ses)
			}
		}
	}
}

//启动服务
func (l *ListenerTcp) Start() error {
	ln, err := net.Listen("tcp", l.addr)
	if err != nil {
		return err
	}

	l.listener = ln
	l.running = true

	//定时处理过期的session
	l.ticker = time.NewTicker(ListenerClearSessionSec * time.Second)
	go l.clearSession()

	for l.running {
		con, err := l.listener.Accept()
		if err != nil {
			continue
		}

		//绑定session
		l.fdcounter++
		ses := NewSession(l.fdcounter, con)
		ses.CallBack = l.CallBack
		ses.OnClose = func(ises common.ISession) {
			l.mutex.Lock()
			delete(l.sessions, ises.ID())
			l.mutex.Unlock()
		}
		l.mutex.Lock()
		l.sessions[l.fdcounter] = ses
		l.mutex.Unlock()
	}
	return nil
}

func NewListener(address string) *ListenerTcp {
	self := &ListenerTcp{
		addr:      address,
		fdcounter: 0,
		log:       utils.NewLogger(),
		running:   false,
		sessions:  make(map[uint64]*Session),
	}
	return self
}

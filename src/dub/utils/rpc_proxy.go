package utils

import (
	"net/rpc"
	"sync"
	"errors"
	"time"
)

//rpc代理
type RpcProxy struct {
	net     string       //协议类型
	address string       //地址
	client  *rpc.Client  //客户端
	mutex   sync.RWMutex //重连锁
	iskeep  bool         //是否连接
}

func (r *RpcProxy) Call(serviceMethod string, args interface{}, reply interface{}) error {
	if !r.iskeep {
		return errors.New("rpc proxy not started:" + serviceMethod)
	}

	r.mutex.RLock()
	err := r.client.Call(serviceMethod, args, reply)
	if err != nil && err.Error() == "connection is shut down" {
		r.mutex.RUnlock()
		r.mutex.Lock()
		r.client.Close()
		ch := make(chan bool)
		go func(ch chan bool) {
			for r.iskeep {
				client, err := rpc.DialHTTP(r.net, r.address)
				if err == nil {
					r.client = client
					r.mutex.Unlock()
					ch <- true
					close(ch)
					break
				}
			}
		}(ch)

		//连接成功之后，重新执行。同时支持超时
		select {
		case <-time.After(time.Second * 3):

			return err
		case <-ch:
			r.mutex.RLock()
			err = r.client.Call(serviceMethod, args, reply)
			r.mutex.RUnlock()
			break
		}
	} else {
		r.mutex.RUnlock()
	}
	return err
}

func (r *RpcProxy) Status() bool {
	return r.iskeep
}

func (r *RpcProxy) Start() error {
	client, err := rpc.DialHTTP(r.net, r.address)
	if err == nil {
		r.client = client
		r.iskeep = true
		return nil
	}
	return err
}

var _instance *RpcProxy

func NewRpcProxy(network, address string) *RpcProxy {
	if _instance == nil {
		_instance = &RpcProxy{
			net:     network,
			address: address,
			iskeep:  false,
		}
	}
	return _instance
}

package dbrpc

import (
	"dub/utils"
)

//与代理有关的数据库操作
type AgentRpc struct {
	log *utils.Logger
}

func (a *AgentRpc) Find(args *uint8, reply *uint8) error {
	return nil
}

var d_agentRpc *AgentRpc

func NewAgentRpc() *AgentRpc {
	if d_agentRpc == nil {
		d_agentRpc = &AgentRpc{
			log: utils.NewLogger(),
		}
	}
	return d_agentRpc
}

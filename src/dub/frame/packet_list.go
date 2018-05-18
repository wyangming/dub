package frame

import (
	"sync"
)

type PacketList struct {
	list      []*Packet
	listGuard sync.Mutex
	listCond  *sync.Cond
	isDone    bool
}

func (this *PacketList) Add(p *Packet) {
	if !this.isDone {
		this.listGuard.Lock()
		this.list = append(this.list, p)
		this.listGuard.Unlock()

		this.listCond.Signal()
	}

}

func (this *PacketList) Reset() {
	this.list = this.list[0:0]
}

func (this *PacketList) BeginPick() []*Packet {
	this.listGuard.Lock()

	for len(this.list) == 0 {
		this.listCond.Wait()
	}

	this.listGuard.Unlock()
	this.listGuard.Lock()

	return this.list
}

func (this *PacketList) EndPick() {
	this.Reset()
	this.listGuard.Unlock()
}

func (this *PacketList) Close() {
	p := MakePacket(InvliadPacketMainID, 0, nil)
	this.listGuard.Lock()
	this.list = append(this.list, p)
	this.listGuard.Unlock()

	this.isDone = true
	this.listCond.Signal()
}

func NewPacketList() *PacketList {
	this := &PacketList{
		isDone: false,
	}
	this.listCond = sync.NewCond(&this.listGuard)

	return this
}

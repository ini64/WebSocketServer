package Packet

import (
	"sync"
)

//HeaderPool HeaderPool
var HeaderPool = sync.Pool{
	// New is called when a new instance is needed
	New: func() interface{} {
		return new(Header)
	},
}

// GetHeader GetHeader
func GetHeader() *Header {
	return HeaderPool.Get().(*Header)
}

//Release Release
func (v *Header) Release() {
	HeaderPool.Put(v)
}

//CSEnterPool CSEnterPool
var CSEnterPool = sync.Pool{
	// New is called when a new instance is needed
	New: func() interface{} {
		return new(CS_Enter)
	},
}

// GetCSEnter fetches a buffer from the pool
func GetCSEnter() *CS_Enter {
	return CSEnterPool.Get().(*CS_Enter)
}

//Release Release
func (v *CS_Enter) Release() {
	CSEnterPool.Put(v)
}

//CSEnterAckPool CSEnterAckPool
var CSEnterAckPool = sync.Pool{
	// New is called when a new instance is needed
	New: func() interface{} {
		return new(CS_Enter_Ack)
	},
}

// GetCSEnterAck GetCSEnterAck
func GetCSEnterAck() *CS_Enter_Ack {
	return CSEnterAckPool.Get().(*CS_Enter_Ack)
}

//Release Release
func (v *CS_Enter_Ack) Release() {
	CSEnterAckPool.Put(v)
}

//CSBroadcastPool CSBroadcastPool
var CSBroadcastPool = sync.Pool{
	// New is called when a new instance is needed
	New: func() interface{} {
		return new(CS_Broadcast)
	},
}

// GetCSBroadcast GetCSBroadcast
func GetCSBroadcast() *CS_Broadcast {
	return CSBroadcastPool.Get().(*CS_Broadcast)
}

//Release Release
func (v *CS_Broadcast) Release() {
	CSBroadcastPool.Put(v)
}

//CSBroadcastAckPool CSBroadcastAckPool
var CSBroadcastAckPool = sync.Pool{
	// New is called when a new instance is needed
	New: func() interface{} {
		return new(CS_Broadcast_Ack)
	},
}

// GetCSBroadcastAck GetCSBroadcastAck
func GetCSBroadcastAck() *CS_Broadcast_Ack {
	return CSBroadcastAckPool.Get().(*CS_Broadcast_Ack)
}

//Release Release
func (v *CS_Broadcast_Ack) Release() {
	CSBroadcastAckPool.Put(v)
}

//SCBroadcastPool SCBroadcastPool
var SCBroadcastPool = sync.Pool{
	// New is called when a new instance is needed
	New: func() interface{} {
		return new(SC_Broadcast)
	},
}

// GetSCBroadcast GetSCBroadcast
func GetSCBroadcast() *SC_Broadcast {
	return SCBroadcastPool.Get().(*SC_Broadcast)
}

//Release Release
func (v *SC_Broadcast) Release() {
	SCBroadcastPool.Put(v)
}

//SCEnterSecondPool SCEnterSecondPool
var SCEnterSecondPool = sync.Pool{
	// New is called when a new instance is needed
	New: func() interface{} {
		return new(SC_Enter_Second)
	},
}

// GetSCEnterSecond GetSCEnterSecond
func GetSCEnterSecond() *SC_Enter_Second {
	return SCEnterSecondPool.Get().(*SC_Enter_Second)
}

//Release Release
func (v *SC_Enter_Second) Release() {
	SCEnterSecondPool.Put(v)
}

//Release 메모리 반환
func Release(v interface{}) {
	switch packet := v.(type) {
	case *Header:
		packet.Release()

	case *CS_Enter:
		packet.Release()
	case *CS_Enter_Ack:
		packet.Release()
	case *SC_Enter_Second:
		packet.Release()

	case *CS_Broadcast:
		packet.Release()
	case *CS_Broadcast_Ack:
		packet.Release()
	case *SC_Broadcast:
		packet.Release()

	}
}

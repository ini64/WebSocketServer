package lib

import (
	"sync"
)

//Broadcast 브로드 케스트 메세지
type Broadcast struct {
	Client  *Client
	Message string
}

// BroadcastPool sync.Pool
var BroadcastPool = sync.Pool{
	New: func() interface{} {
		return new(Broadcast)
	},
}

// GetBroadcast fetches a buffer from the pool
func GetBroadcast() *Broadcast {
	return BroadcastPool.Get().(*Broadcast)
}

//Release Release
func (m *Broadcast) Release() {
	BroadcastPool.Put(m)
}

//SessionEnter 세션에 들어 가는 패킷
type SessionEnter struct {
	Client *Client
}

// SessionEnterPool sync.Pool
var SessionEnterPool = sync.Pool{
	New: func() interface{} {
		return new(SessionEnter)
	},
}

// GetSessionEnter fetches a buffer from the pool
func GetSessionEnter() *SessionEnter {
	return SessionEnterPool.Get().(*SessionEnter)
}

//Release Release
func (m *SessionEnter) Release() {
	SessionEnterPool.Put(m)
}

//SessionLeave 세션에 들어 가는 패킷
type SessionLeave struct {
	Client *Client
}

// SessionLeavePool sync.Pool
var SessionLeavePool = sync.Pool{
	New: func() interface{} {
		return new(SessionLeave)
	},
}

// GetSessionLeave fetches a buffer from the pool
func GetSessionLeave() *SessionLeave {
	return SessionLeavePool.Get().(*SessionLeave)
}

//Release Release
func (m *SessionLeave) Release() {
	SessionLeavePool.Put(m)
}

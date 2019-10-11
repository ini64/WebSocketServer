package lib

import (
	"net/http"
	"sync"
)

//EndPoint 전체 정보
type EndPoint struct {
	Conf
	ClientWait   sync.WaitGroup
	ListenerWait sync.WaitGroup

	WSListenerMake      chan bool
	SessionListenerMake chan bool

	Server *http.Server

	SessionChannel chan interface{}

	RecvCount  int64
	SendCount  int64
	TotalCount int64

	slowly chan bool
}

//NewEndPoint 생성
func NewEndPoint() *EndPoint {

	endPoint := &EndPoint{
		WSListenerMake:      make(chan bool),
		SessionListenerMake: make(chan bool),

		SessionChannel: make(chan interface{}, 2048),

		Conf: Conf{WSBind: "7000"},
	}

	return endPoint
}

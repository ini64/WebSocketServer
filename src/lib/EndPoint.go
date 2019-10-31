package lib

import (
	"fmt"
	"net/http"
	"sync"
)

//EndPoint 전체 정보
type EndPoint struct {
	*Conf
	*Redis
	ClientWait   sync.WaitGroup
	ListenerWait sync.WaitGroup

	WSListenerMake      chan bool
	SessionListenerMake chan bool
	RedisListenerMake   chan bool

	Server *http.Server

	SessionChannel chan interface{}

	RecvCount  int64
	SendCount  int64
	TotalCount int64

	ServerID int64
}

//NewEndPoint 생성
func NewEndPoint() *EndPoint {
	conf := &Conf{}
	if conf.Load("Conf.json") == false {
		fmt.Println("Not Found Conf.json")
		return nil
	}

	endPoint := &EndPoint{
		WSListenerMake:      make(chan bool),
		SessionListenerMake: make(chan bool),
		RedisListenerMake:   make(chan bool),
		SessionChannel:      make(chan interface{}, 2048),
		Conf:                conf,
		Redis:               conf.GetRedis(),
	}

	return endPoint
}

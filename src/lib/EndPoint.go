package lib

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/go-redis/redis"
)

//EndPoint 전체 정보
type EndPoint struct {
	*Conf
	RedisClient  *redis.Client
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
	conf := &Conf{}
	if conf.Load("Conf.json") == false {
		fmt.Println("Not Found Conf.json")
		return nil
	}

	client := redis.NewClient(&redis.Options{
		Addr:     conf.RedisAddr,
		Password: conf.RedisPassword, // no password set
		DB:       0,                  // use default DB
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)

	endPoint := &EndPoint{
		WSListenerMake:      make(chan bool),
		SessionListenerMake: make(chan bool),
		SessionChannel:      make(chan interface{}, 2048),
		Conf:                conf,
		RedisClient:         client,
	}

	return endPoint
}

package lib

import (
	"Packet"
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"nhooyr.io/websocket"
)

//Client 클라이언트 접근 정보
type Client struct {
	//서버 관리 객체
	EndPoint *EndPoint

	//게임 UID
	GameUID uint64

	//네트워크 연결
	Conn *websocket.Conn

	//전송자와 싱크를 맞추기 위해
	SyncSender chan bool

	//세션과 싱크를 맞추기 위해
	SyncSession chan bool

	SendListSize int

	//소켓으로 보낼 리스트
	SendList chan interface{}

	//전송자를 종료시키기 위한 객체
	Context context.Context

	//전송자 종료 함수
	Cancel context.CancelFunc

	//전송 가능 한지?
	Transferable int32

	//sync.RWMutex
}

// ClientPool sync.Pool
var ClientPool = sync.Pool{
	New: func() interface{} {
		client := new(Client)
		client.Make()
		return client
	},
}

//Make 처음 생성시
func (c *Client) Make() {
	c.SyncSender = make(chan bool)
	c.SyncSession = make(chan bool)
	c.SendListSize = 8
	c.SendList = make(chan interface{}, c.SendListSize)
}

//Reset Reset
func (c *Client) Reset() {

	c.EndPoint = nil
	c.GameUID = 0
	c.Transferable = 0

	c.SyncSender = make(chan bool)
	c.SyncSession = make(chan bool)

	//닫고 해야 Lock이 안걸림
	close(c.SendList)
	for packet := range c.SendList {
		if packet != nil {
			Packet.Release(packet)
		}
	}
	c.SendList = make(chan interface{}, c.SendListSize)

	c.Context = nil
	c.Cancel = nil
	c.Conn = nil
}

//Send 만약에 최대치에 도달하면 실패를 리턴한다
func (c *Client) Send(v interface{}) bool {

	//전송자가 닫힘
	transferable := atomic.LoadInt32(&c.Transferable)
	if transferable == 0 {
		Packet.Release(v)
		return false
	}

	//최대치 도달
	if len(c.SendList) == c.SendListSize {
		Packet.Release(v)
		return false
	}

	c.SendList <- v
	return true
}

//NetClose 세션에서 닫던가 세션과 연결이 끊어지면 호출
// func (c *Client) NetClose() {
// 	if c.Net != nil {
// 		c.Net.Close()
// 		c.Net = nil
// 	}
// }

//Context Context 얻기
// func (c *Client) Context() context.Context {
// 	return c.Net.Request.Context()
// }

// NewClient 개인 유저 처리
func (e *EndPoint) NewClient(conn *websocket.Conn, slowy chan bool) {

	var (
		gameUID uint64
		seq     uint32
	)

	client := ClientPool.Get().(*Client)
	client.EndPoint = e
	client.Conn = conn

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("NewClient", "RunTime Panic", string(Stack()), err)
		}
		//<-slowy

		// if gameUID == 0 {
		// 	fmt.Println("test")
		// }

		sessionLeave := GetSessionLeave()
		sessionLeave.Client = client
		e.SessionChannel <- sessionLeave

		for range client.SyncSession {
		}
		for range client.SyncSender {
		}

		//자원 정리
		client.Reset()

		//반납
		ClientPool.Put(client)
		e.ClientWait.Done()

		//fmt.Println("exit", gameUID)
	}()

	//전송자 생성
	client.Context, client.Cancel = context.WithCancel(context.Background())

	//타임 아웃 설정
	firstCtx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	v, err := ReadWS(firstCtx, client.Conn, &e.RecvCount)
	if err != nil {
		close(client.SyncSender)
		return
	}

	switch message := v.(type) {
	case *Packet.CS_Enter:
		gameUID = uint64(message.GameUID)
		seq = message.SEQ
	}
	Packet.Release(v)

	if gameUID == 0 {
		close(client.SyncSender)
		return
	}

	Packet.Log(gameUID, "CS_Enter", v)

	dbSeq, result := e.Redis.GetSEQ(gameUID)
	if result == false {
		close(client.SyncSender)
		return
	}

	if dbSeq != seq {
		close(client.SyncSender)
		return
	}

	client.GameUID = gameUID

	go client.SendWorker(e)
	<-client.SyncSender

	//세션에 등록
	sessionEnter := GetSessionEnter()
	sessionEnter.Client = client
	e.SessionChannel <- sessionEnter
	<-client.SyncSession

	for {
		select {
		case <-client.Context.Done():
			return
		default:
			//5분동안 반응 없으면 에러처리
			ctx, mainCancel := context.WithTimeout(context.Background(), time.Second*600)
			defer mainCancel()

			v, err := ReadWS(ctx, client.Conn, &e.RecvCount)
			if err != nil {
				return
			}

			switch message := v.(type) {
			case *Packet.CS_Broadcast:

				Packet.Log(gameUID, "CS_Broadcast", v)

				broadcast := GetBroadcast()
				broadcast.Client = client
				broadcast.Message = message.Message
				e.SessionChannel <- broadcast
			}
			Packet.Release(v)
		}
	}
}

//SendWorker 전송용
func (c *Client) SendWorker(e *EndPoint) {

	defer func() {
		atomic.StoreInt32(&c.Transferable, 0)
		close(c.SyncSender)
		//fmt.Println("Exit SendWorker", gameUID)
	}()

	conn := c.Conn
	gameUID := c.GameUID

	atomic.StoreInt32(&c.Transferable, 1)
	c.SyncSender <- true

	for {
		select {
		case v := <-c.SendList:

			ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
			defer cancel()

			err := WriteWS(ctx, conn, v, &e.SendCount, gameUID)

			//보낼 수 없으면 종료
			if err != nil {
				return
			}
		case <-c.Context.Done():
			return
		}
	}
}

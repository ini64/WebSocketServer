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
	SyncSenderStart chan bool

	//전송자와 싱크를 맞추기 위해
	SyncSenderEnd chan bool

	//세션과 싱크를 맞추기 위해
	SyncSessionStart chan bool

	//세션과 싱크를 맞추기 위해
	SyncSessionEnd chan bool

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
	c.SyncSenderStart = make(chan bool, 2)
	c.SyncSenderEnd = make(chan bool, 2)
	c.SyncSessionStart = make(chan bool, 2)
	c.SyncSessionEnd = make(chan bool, 2)
	c.SendListSize = 8
	c.SendList = make(chan interface{}, c.SendListSize)
}

//Reset Reset
func (c *Client) Reset() {

	c.EndPoint = nil
	c.GameUID = 0
	c.Transferable = 0

	close(c.SyncSenderStart)
	close(c.SyncSenderEnd)
	close(c.SyncSessionStart)
	close(c.SyncSessionStart)

	c.SyncSenderStart = make(chan bool, 2)
	c.SyncSenderEnd = make(chan bool, 2)
	c.SyncSessionStart = make(chan bool, 2)
	c.SyncSessionEnd = make(chan bool, 2)

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
		//Token   string
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

		<-client.SyncSessionEnd
		<-client.SyncSenderEnd

		//자원 정리
		client.Reset()

		//반납
		ClientPool.Put(client)
		e.ClientWait.Done()

		//fmt.Println("exit", gameUID)
	}()

	//전송자 생성
	client.Context, client.Cancel = context.WithCancel(context.Background())

	go client.SendWorker(e)
	<-client.SyncSenderStart

	//타임 아웃 설정
	firstCtx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	v, err := ReadWS(firstCtx, client.Conn, &e.RecvCount)
	if err != nil {
		//fmt.Println("invalid packet")
		return
	}

	switch message := v.(type) {
	case *Packet.CS_Enter:
		gameUID = uint64(message.GameUID)
		//Token = message.Token
	}
	Packet.Release(v)

	if gameUID == 0 {
		return
	}

	client.GameUID = gameUID

	//todo 여기서 보내온 데이터의 적합성을 체크해야 하지만 나중에 하자. ase256으로 보낸 gameUID와 Token을 비교
	// if len(Token) != 0 {
	// 	return
	// }

	//세션에 등록
	sessionEnter := GetSessionEnter()
	sessionEnter.Client = client
	e.SessionChannel <- sessionEnter
	<-client.SyncSessionStart

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
		c.SyncSenderEnd <- true
		//fmt.Println("Exit SendWorker", gameUID)
	}()

	conn := c.Conn

	atomic.StoreInt32(&c.Transferable, 1)
	c.SyncSenderStart <- true

	for {
		select {
		case v := <-c.SendList:
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
			defer cancel()

			err := WriteWS(ctx, conn, v, &e.SendCount)

			//보낼 수 없으면 종료
			if err != nil {
				return
			}
		case <-c.Context.Done():
			return
		}
	}
}

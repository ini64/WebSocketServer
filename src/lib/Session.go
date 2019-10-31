package lib

import (
	"Packet"
	"context"
	"fmt"
	"sync/atomic"
)

//SessionListener 작업
func (e *EndPoint) SessionListener(ctx context.Context) {
	defer func() {
		e.ListenerWait.Done()
	}()

	clients := make(map[uint64]*Client)

	e.SessionListenerMake <- true

	for {
		if e.SessionParser(ctx, clients) {
			return
		}
	}
}

//SessionParser 루프 터지면 망이라 recover처리 추가
func (e *EndPoint) SessionParser(ctx context.Context, clients map[uint64]*Client) (exit bool) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("SessionParser", "RunTime Panic", string(Stack()), err)
		}
	}()

	exit = false

	select {
	case recvMessage := <-e.SessionChannel:
		switch message := recvMessage.(type) {

		case *SessionEnter:
			defer message.Release()

			gameUID := message.Client.GameUID
			client, ok := clients[gameUID]
			if ok {
				//클라가 소켓 끊을 수 있도록 패킷 전송
				scEnterSecond := Packet.GetSCEnterSecond()
				client.Send(scEnterSecond)
				delete(clients, gameUID)
			}

			//다른 서버에 있는 채널에게 세션을 끊으라고 전달한다.
			err := e.Redis.Monitor.Publish("session", fmt.Sprintf("%d|%d", e.ServerID, gameUID)).Err()
			if err != nil {
				scEnterSecond := Packet.GetSCEnterSecond()
				message.Client.Send(scEnterSecond)
			} else {
				//세션에 등록
				clients[gameUID] = message.Client

				//세션 등록 완료 메세지 보내기
				csEnterAck := Packet.GetCSEnterAck()
				message.Client.Send(csEnterAck)

				//대기 중인 스레드에게 시그널
				message.Client.SyncSession <- true
				atomic.StoreInt64(&e.TotalCount, int64(len(clients)))
			}

		case *SessionLeave:
			defer message.Release()

			gameUID := message.Client.GameUID
			client, ok := clients[gameUID]
			if ok {
				//등록된 클라이언트가 종료를 요청한다면
				if client == message.Client {
					delete(clients, gameUID)
				}
			}

			//스레드 종료하라고 시그널
			message.Client.Cancel()
			//대기중인 스레드에게 시그널
			close(message.Client.SyncSession)

			atomic.StoreInt64(&e.TotalCount, int64(len(clients)))

		case *Broadcast:
			defer message.Release()

			gameUID := message.Client.GameUID
			count := int(0)

			client, ok := clients[gameUID]
			if ok {
				if client == message.Client {
					//나 외의 다른 클라이언트 찾기
					for otherGameUID, otherClient := range clients {
						if gameUID != otherGameUID {

							scBroadcast := Packet.GetSCBroadcast()
							scBroadcast.Message = message.Message
							otherClient.Send(scBroadcast)
							count++
						}
					}
					//완료 시그널 보내기
					csBroadcastAck := Packet.GetCSBroadcastAck()
					client.Send(csBroadcastAck)
				}
			}
			atomic.StoreInt64(&e.TotalCount, int64(len(clients)))

		case *ServerEnter:
			defer message.Release()

			gameUID := message.GameUID
			client, ok := clients[gameUID]
			if ok {
				//클라가 소켓 끊을 수 있도록 패킷 전송
				scEnterSecond := Packet.GetSCEnterSecond()
				client.Send(scEnterSecond)

				delete(clients, gameUID)
			}
		}

	case <-ctx.Done():
		exit = true
	}

	return exit
}

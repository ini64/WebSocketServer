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
	case ss := <-e.SessionChannel:
		switch x := ss.(type) {

		case *SessionEnter:
			gameUID := x.Client.GameUID
			client, ok := clients[gameUID]
			if ok {
				client.Cancel()
				client.SyncSessionEnd <- true
			}
			clients[gameUID] = x.Client

			csEnterAck := Packet.GetCSEnterAck()
			x.Client.Send(csEnterAck)

			x.Client.SyncSessionStart <- true

			atomic.StoreInt64(&e.TotalCount, int64(len(clients)))

			x.Release()

		case *SessionLeave:
			gameUID := x.Client.GameUID
			client, ok := clients[gameUID]
			if ok && client == x.Client {
				delete(clients, gameUID)
				client.Cancel()
				client.SyncSessionEnd <- true
			} else {
				x.Client.Cancel()
				x.Client.SyncSessionEnd <- true
			}
			atomic.StoreInt64(&e.TotalCount, int64(len(clients)))

			x.Release()

		case *Broadcast:
			gameUID := x.Client.GameUID
			count := int(0)

			client, ok := clients[gameUID]
			if ok {
				if client == x.Client {
					for otherGameUID, client := range clients {
						if gameUID != otherGameUID {

							scBroadcast := Packet.GetSCBroadcast()
							scBroadcast.Message = x.Message
							client.Send(scBroadcast)
							count++
						}
						if count > 3 {
							break
						}
					}

					csBroadcastAck := Packet.GetCSBroadcastAck()
					client.Send(csBroadcastAck)
				}
			}
			atomic.StoreInt64(&e.TotalCount, int64(len(clients)))

			x.Release()
		}

	case <-ctx.Done():
		exit = true
	}

	return exit
}

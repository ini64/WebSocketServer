package client

import (
	"Packet"
	"context"
	"fmt"
	"lib"
	"math/rand"
	"net"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"nhooyr.io/websocket"
)

func PerfomanceTest() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("run tim panic", err, string(lib.Stack()))
		}
	}()

	var wg sync.WaitGroup

	total := uint64(10000)
	run := int(1024)
	slowly := int(256)

	totalRun := make(chan bool, run)
	sameTime := make(chan bool, slowly)
	var gameUID uint64

	for i := uint64(1); i < (total + 1); i++ {
		wg.Add(1)
		totalRun <- true
		sameTime <- true

		//time.Sleep(time.Nanosecond)
		go Client(&wg, &gameUID, totalRun, sameTime)
	}
	wg.Wait()
}

// var timeout = time.Duration(60 * time.Second)

// func dialTimeout(network, addr string) (net.Conn, error) {
// 	return net.DialTimeout(network, addr, timeout)
// }

// Write writes the protobuf message v to c.
// It will reuse buffers to avoid allocations.

func Client(wg *sync.WaitGroup, gameUID *uint64, totalRun chan bool, sameTime chan bool) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("run tim panic", err, string(lib.Stack()))
		}
		wg.Done()
		<-totalRun
	}()

	tr := &http.Transport{
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSHandshakeTimeout:   30 * time.Second,
		ResponseHeaderTimeout: 30 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		//MaxIdleConns:          7,
	}

	options := &websocket.DialOptions{
		HTTPClient: &http.Client{
			Transport: tr,
		},
	}
	ctx, cancel1 := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel1()
	c, _, err := websocket.Dial(ctx, "ws://0.0.0.0:7000", options)
	if err != nil {
		fmt.Println("websocket.Dial", err)
		<-sameTime
		return
	}
	defer c.Close(websocket.StatusInternalError, "the sky is falling")

	<-sameTime

	csEnter := Packet.GetCSEnter()
	csEnter.GameUID = uint64(atomic.AddUint64(gameUID, 1))

	ctx, cancel2 := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel2()
	err = lib.WriteWS(ctx, c, csEnter, nil)
	if err != nil {
		fmt.Println("csEnter Error", err)
		return
	}

	ctx, cancel3 := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel3()
	v, err := lib.ReadWS(ctx, c, nil)

	switch v.(type) {
	case *Packet.CS_Enter_Ack:
	default:
		fmt.Println("client error 1", err)
		return
	}
	Packet.Release(v)

	for i := 0; i < 10; i++ {

		csBroadcast := Packet.GetCSBroadcast()
		csBroadcast.Message = fmt.Sprintf("%s%d", "hi", int(*gameUID))

		ctx, cancel4 := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel4()
		err = lib.WriteWS(ctx, c, csBroadcast, nil)
		if err != nil {
			fmt.Println("CSBoradcast Error", err)
			return
		}

		ctx, cancel5 := context.WithTimeout(context.Background(), time.Second*120)
		defer cancel5()
		v, err := lib.ReadWS(ctx, c, nil)
		if err != nil {
			fmt.Println("client error 2", err)
			return
		}
		defer Packet.Release(v)

		switch v.(type) {
		case *Packet.SC_Broadcast:
		case *Packet.CS_Broadcast_Ack:
		default:
			fmt.Println("Not Recv scBoradcast")
			return
		}
		sleepTime := time.Duration(int(5 * rand.Float64()))
		time.Sleep(time.Second * sleepTime)
	}

	// message2 := lib.GetMessage()
	// defer message2.Release()
	// err = wspb.Read(ctx, c, message2)
	// if err != nil {
	// 	fmt.Println("message2 Error", err)
	// 	return
	// }

	// switch x := message2.Avatar.(type) {
	// case *Packet.Message_SC_Broadcast:
	// 	if x.SC_Broadcast != nil {
	// 	}
	// 	break
	// case *Packet.Message_CS_Broadcast_Ack:
	// 	if x.CS_Broadcast_Ack != nil {
	// 	}
	// 	break
	// default:
	// 	fmt.Println("Not Recv scBoradcast")
	// 	return
	// }

	c.Close(websocket.StatusNormalClosure, "")

}

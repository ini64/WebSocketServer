package main

import (
	"client"
	"context"
	"fmt"
	"lib"
	"os"
	"runtime"
	"sync/atomic"
	"time"
)

func main() {
	endPoint := lib.NewEndPoint()

	runType := os.Args[1]
	switch runType {
	case "server":
		runtime.GOMAXPROCS(2)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		endPoint.ListenerWait.Add(2)

		go endPoint.WSListener()
		go endPoint.SessionListener(ctx)

		<-endPoint.WSListenerMake
		<-endPoint.SessionListenerMake

		go Monitor(endPoint)
		endPoint.ListenerWait.Wait()

	case "client":
		runtime.GOMAXPROCS(runtime.NumCPU())
		client.PerfomanceTest()
		endPoint.ClientWait.Wait()

	case "all":
		runtime.GOMAXPROCS(runtime.NumCPU())

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		endPoint.ListenerWait.Add(2)

		go endPoint.WSListener()
		go endPoint.SessionListener(ctx)

		<-endPoint.WSListenerMake
		<-endPoint.SessionListenerMake

		go Monitor(endPoint)
		client.PerfomanceTest()
		endPoint.ClientWait.Wait()
		endPoint.ListenerWait.Wait()

	}

}

func Monitor(e *lib.EndPoint) {

	var Send, Recv int64
	var m runtime.MemStats

	for {
		select {
		case <-time.After(1 * time.Second):

			recvCount := atomic.LoadInt64(&e.RecvCount)
			sendCount := atomic.LoadInt64(&e.SendCount)
			total := atomic.LoadInt64(&e.TotalCount)

			recv := recvCount - Recv
			Recv = recvCount

			send := sendCount - Send
			Send = sendCount

			t := time.Now()

			runtime.ReadMemStats(&m)

			fmt.Fprintln(os.Stdout, t.Format(time.RFC3339), ",", runtime.NumGoroutine(), ",", recv, ",", send, ",", total)
		}
	}
}

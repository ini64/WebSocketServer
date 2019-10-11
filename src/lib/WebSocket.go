package lib

import (
	"Packet"
	"bytes"
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wspb"
)

// bpool sync.Pool
var bpool = sync.Pool{
	New: func() interface{} {
		return new(Broadcast)
	},
}

// Get returns a buffer from the pool or creates a new one if
// the pool is empty.
func Get() *bytes.Buffer {
	b, ok := bpool.Get().(*bytes.Buffer)
	if !ok {
		b = &bytes.Buffer{}
	}
	return b
}

// Put returns a buffer into the pool.
func Put(b *bytes.Buffer) {
	b.Reset()
	bpool.Put(b)
}

//Net 관리용으로 함수를 만들기 위해
type Net struct {
	Request *http.Request
	Conn    *websocket.Conn
}

//Close 네트워크 닫기
func (n *Net) Close() {
	n.Conn.Close(websocket.StatusNormalClosure, "Call Close")
}

type LimitedListener struct {
	//sync.Mutex
	net.Listener
	sem chan bool
}

func NewLimitedListener(sem chan bool, l net.Listener) *LimitedListener {
	return &LimitedListener{
		Listener: l,
		sem:      sem,
	}
}

// type Listener interface {
// 	// Accept waits for and returns the next connection to the listener.
// 	Accept() (Conn, error)

// 	// Close closes the listener.
// 	// Any blocked Accept operations will be unblocked and return errors.
// 	Close() error

// 	// Addr returns the listener's network address.
// 	Addr() Addr
// }

func (l *LimitedListener) Addr() net.Addr {
	return l.Listener.Addr()
}
func (l *LimitedListener) Close() error {
	return l.Listener.Close()
}
func (l *LimitedListener) Accept() (net.Conn, error) {
	//l.sem <- true // acquire
	return l.Listener.Accept()
}

//WSListener 메세지 대기
func (e *EndPoint) WSListener() {
	defer func() {
		e.ListenerWait.Done()
	}()

	l, err := net.Listen("tcp", "0.0.0.0:"+e.Conf.WSBind)
	if err != nil {
		fmt.Println("failed to listen:", err)
	}

	fmt.Println("TCP Listen:0.0.0.0", e.Conf.WSBind)
	defer l.Close()

	slowy := make(chan bool, 1)
	limitedListener := NewLimitedListener(slowy, l)

	e.Server = &http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			//fmt.Println("serving", r.RemoteAddr)

			c, err := websocket.Accept(w, r, nil)
			if err != nil {
				return
			}
			defer c.Close(websocket.StatusInternalError, "the sky is falling")

			//클라 완전 종료 체크용
			e.ClientWait.Add(1)

			//접속 처리
			e.NewClient(c, slowy)
		}),
		ReadTimeout:  time.Second * 30,
		WriteTimeout: time.Second * 30,
	}
	defer e.Server.Close()

	e.WSListenerMake <- true

	err = e.Server.Serve(limitedListener)
	if err != http.ErrServerClosed {
		fmt.Println("failed to listen and serve:", err)
	}
}

// func Write(ctx context.Context, c *websocket.Conn, v proto.Message) error {
// 	err := write(ctx, c, v)
// 	if err != nil {
// 		return fmt.Errorf("failed to write protobuf: %w", err)
// 	}
// 	return nil
// }

// func write(ctx context.Context, c *websocket.Conn, v proto.Message) error {
// 	// b := Get()
// 	// pb := proto.NewBuffer(b.Bytes())
// 	// defer func() {
// 	// 	Put(b)
// 	// }()
// 	fmt.Printf("write %v \n", v)
// 	out, err := proto.Marshal(v)
// 	if err != nil {
// 		return fmt.Errorf("failed to marshal protobuf: %w", err)
// 	}
// 	//buf := pb.Bytes()
// 	fmt.Println("pb write", hex.EncodeToString(out))
// 	return c.Write(ctx, websocket.MessageBinary, out)
// }

//WriteWS 데이터 쓰기
func WriteWS(ctx context.Context, c *websocket.Conn, v interface{}, counter *int64) error {
	if counter != nil {
		defer atomic.AddInt64(counter, 2)
	}

	defer Packet.Release(v)

	header := Packet.GetHeader()
	defer Packet.Release(header)

	switch packet := v.(type) {
	case *Packet.CS_Enter:
		header.Type = Packet.Header_CS_Enter
		err := wspb.Write(ctx, c, header)
		if err != nil {
			return err
		}
		return wspb.Write(ctx, c, packet)

	case *Packet.CS_Enter_Ack:
		header.Type = Packet.Header_CS_Enter_Ack
		err := wspb.Write(ctx, c, header)
		if err != nil {
			return err
		}
		return wspb.Write(ctx, c, packet)

	case *Packet.CS_Broadcast:
		header.Type = Packet.Header_CS_Broadcast
		err := wspb.Write(ctx, c, header)
		if err != nil {
			return err
		}
		return wspb.Write(ctx, c, packet)

	case *Packet.CS_Broadcast_Ack:
		header.Type = Packet.Header_CS_Broadcast_Ack
		err := wspb.Write(ctx, c, header)
		if err != nil {
			return err
		}
		return wspb.Write(ctx, c, packet)

	case *Packet.SC_Broadcast:
		header.Type = Packet.Header_SC_Broadcast
		err := wspb.Write(ctx, c, header)
		if err != nil {
			return err
		}
		return wspb.Write(ctx, c, packet)
	}

	return errors.New("not support type")
}

// func Read(ctx context.Context, c *websocket.Conn, v proto.Message) error {
// 	err := read(ctx, c, v)
// 	if err != nil {
// 		return fmt.Errorf("failed to read protobuf: %w", err)
// 	}
// 	return nil
// }

// func read(ctx context.Context, c *websocket.Conn, v proto.Message) error {
// 	typ, r, err := c.Reader(ctx)
// 	if err != nil {
// 		return err
// 	}

// 	if typ != websocket.MessageBinary {
// 		c.Close(websocket.StatusUnsupportedData, "can only accept binary messages")
// 		return fmt.Errorf("unexpected frame type for protobuf (expected %v): %v", websocket.MessageBinary, typ)
// 	}

// 	b := Get()
// 	defer Put(b)

// 	_, err = b.ReadFrom(r)
// 	if err != nil {
// 		return err
// 	}
// 	fmt.Println("pb read", hex.EncodeToString(b.Bytes()))
// 	err = proto.Unmarshal(b.Bytes(), v)
// 	if err != nil {
// 		c.Close(websocket.StatusInvalidFramePayloadData, "failed to unmarshal protobuf")
// 		return fmt.Errorf("failed to unmarshal protobuf: %w", err)
// 	}

// 	return nil
// }

//ReadWS 데이터 읽기
func ReadWS(ctx context.Context, c *websocket.Conn, counter *int64) (interface{}, error) {
	if counter != nil {
		defer atomic.AddInt64(counter, 2)
	}

	header := Packet.GetHeader()
	defer Packet.Release(header)

	err := wspb.Read(ctx, c, header)
	if err != nil {
		return nil, err
	}

	switch header.Type {
	case Packet.Header_CS_Enter:
		csEnter := Packet.GetCSEnter()
		err := wspb.Read(ctx, c, csEnter)
		return csEnter, err

	case Packet.Header_CS_Enter_Ack:
		csEnterAck := Packet.GetCSEnterAck()
		err := wspb.Read(ctx, c, csEnterAck)
		return csEnterAck, err

	case Packet.Header_CS_Broadcast:
		csBroadcast := Packet.GetCSBroadcast()
		err := wspb.Read(ctx, c, csBroadcast)
		return csBroadcast, err

	case Packet.Header_CS_Broadcast_Ack:
		csBroadcastAck := Packet.GetCSBroadcastAck()
		err := wspb.Read(ctx, c, csBroadcastAck)
		return csBroadcastAck, err

	case Packet.Header_SC_Broadcast:
		scBroadcast := Packet.GetSCBroadcast()
		err := wspb.Read(ctx, c, scBroadcast)
		return scBroadcast, err
	}
	return nil, errors.New("invalid packet")
}

package lib

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/go-redis/redis"
)

//Redis  내부용 redis 객체
type Redis struct {
	Token   *redis.Client
	Monitor *redis.Client
}

//GetSEQ 시퀀스 번호 얻기
func (r *Redis) GetSEQ(gameUID uint64) (uint32, bool) {
	seqKey := fmt.Sprintf("seq:%d", gameUID)
	val, err := r.Token.Get(seqKey).Result()
	if err != nil {
		return 0, false
	}

	dbSeq, err := strconv.ParseUint(val, 10, 32)
	if err != nil {
		return 0, false
	}
	return uint32(dbSeq), true
}

//RedisListener 자신의 상태 업데이트를 위해
func (e *EndPoint) RedisListener(ctx context.Context) {
	defer func() {
		e.ListenerWait.Done()
	}()

	seqKey := "SessiongServerSEQ"
	serverID, err := e.Redis.Monitor.Incr(seqKey).Result()

	if err != nil {
		panic(err)
	}

	e.ServerID = serverID

	serverKey := fmt.Sprintf("SessionServer:%d", serverID)

	ticker := time.NewTicker(3 * time.Minute)

	info := make(map[string]interface{})

	dt := time.Now()
	recvCount := atomic.LoadInt64(&e.RecvCount)
	sendCount := atomic.LoadInt64(&e.SendCount)
	total := atomic.LoadInt64(&e.TotalCount)

	info["total"] = total
	info["recv"] = recvCount
	info["send"] = sendCount
	info["time"] = dt.String()
	e.Redis.Monitor.HMSet(serverKey, info).Result()
	e.Redis.Monitor.Expire(serverKey, time.Minute*5)

	atomic.StoreInt64(&e.RecvCount, 0)
	atomic.StoreInt64(&e.SendCount, 0)

	channel := e.Redis.Monitor.Subscribe("session").Channel()
	e.RedisReceice(ctx, channel)

	e.RedisListenerMake <- true

	for {
		select {
		case <-ticker.C:
			dt := time.Now()
			recvCount := atomic.LoadInt64(&e.RecvCount)
			sendCount := atomic.LoadInt64(&e.SendCount)
			total := atomic.LoadInt64(&e.TotalCount)

			info["total"] = total
			info["recv"] = recvCount
			info["send"] = sendCount
			info["time"] = dt.String()
			e.Redis.Monitor.HMSet(serverKey, info).Result()
			e.Redis.Monitor.Expire(serverKey, time.Minute*5)

			atomic.StoreInt64(&e.RecvCount, 0)
			atomic.StoreInt64(&e.SendCount, 0)
		case <-ctx.Done():
			return
		}
	}
}

//RedisReceice pub/sub
func (e *EndPoint) RedisReceice(ctx context.Context, channel <-chan *redis.Message) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-channel:

			array := strings.Split(msg.Payload, "|")
			stringServerID := array[0]
			stringGameUID := array[1]

			serverID, err := strconv.ParseInt(stringServerID, 10, 64)
			if err != nil {
				continue
			}
			gameUID, err := strconv.ParseUint(stringGameUID, 10, 64)
			if err != nil {
				continue
			}

			//서버가 다르면 작업 시작
			if serverID != e.ServerID {
				ssEnter := GetServerEnter()
				ssEnter.GameUID = gameUID
				ssEnter.ServerID = serverID

				e.SessionChannel <- ssEnter
			}
		}
	}

}

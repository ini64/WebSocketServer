package lib

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/go-redis/redis"
)

//Conf 설정
type Conf struct {
	WSBind string

	TokenRedisAddr     string
	TokenRedisPassword string
	TokenRedisDB       int

	MonitorRedisAddr     string
	MonitorRedisPassword string
	MonitorRedisDB       int
}

//Load 파일에서 설정값을 읽어 온다.
func (c *Conf) Load(filePath string) bool {
	file, _ := os.Open(filePath)
	decoder := json.NewDecoder(file)
	err := decoder.Decode(c)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true
}

//GetRedis redis 얻기
func (c *Conf) GetRedis() *Redis {
	r := &Redis{}

	r.Token = redis.NewClient(&redis.Options{
		Addr:     c.TokenRedisAddr,
		Password: c.TokenRedisPassword, // no password set
		DB:       c.TokenRedisDB,       // use default DB
	})

	r.Monitor = redis.NewClient(&redis.Options{
		Addr:     c.MonitorRedisAddr,
		Password: c.MonitorRedisPassword, // no password set
		DB:       c.MonitorRedisDB,       // use default DB
	})

	pong, err := r.Token.Ping().Result()
	fmt.Println(pong, err)

	pong, err = r.Monitor.Ping().Result()
	fmt.Println(pong, err)

	return r
}

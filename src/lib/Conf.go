package lib

import (
	"encoding/json"
	"fmt"
	"os"
)

//Conf 설정
type Conf struct {
	WSBind        string
	RedisAddr     string
	RedisPassword string
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

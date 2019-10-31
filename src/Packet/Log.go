package Packet

import (
	"encoding/json"
	fmt "fmt"
)

//Log 찍기
func Log(gameUID uint64, packetName string, v interface{}) {
	j, err := json.Marshal(v)
	if err == nil {
		fmt.Printf("{\"Name\":\"%s\", \"GameUID\":%d, \"Packet\":\"%s\"}\n", packetName, gameUID, string(j))
	}
}

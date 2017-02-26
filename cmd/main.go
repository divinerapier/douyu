package main

import (
	"flag"
	"fmt"

	log "qiniupkg.com/x/log.v7"

	"github.com/DivineRapier/douyu"
)

// ZSMJ 52876
// 天使焦 97376

func main() {

	room := flag.Int64("room", 52876, "room_id")

	flag.Parse()

	dy, err := douyu.OpenDanmu(*room)
	if err != nil {
		log.Error(err)
		return
	}
	fmt.Println(dy)
	fmt.Println()
	fmt.Println()

	dy.JoinGroupRequest(0)
	dy.ReceiveResponse()
	dy.HeartBeat()
	wait()
}

func wait() {
	c := make(chan struct{})
	<-c
}

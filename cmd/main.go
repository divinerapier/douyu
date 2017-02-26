package main

import (
	"fmt"

	log "qiniupkg.com/x/log.v7"

	"github.com/DivineRapier/douyu"
)

// ZSMJ 52876
// 天使焦 97376

func main() {
	dy, err := douyu.OpenDanmu(97376)
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

package main

import (
	"fmt"

	"github.com/DivineRapier/douyu"
)

// ZSMJ 52876

func main() {
	dy, err := douyu.OpenDanmu(52876)
	if err != nil {
		fmt.Println(err)
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

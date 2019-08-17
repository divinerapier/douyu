package main

import (
	"flag"
	"os"

	"github.com/divinerapier/douyu/danmaku"
	"github.com/sirupsen/logrus"
)

// ZSMJ 52876
// 天使焦 97376

func main() {
	room := flag.Int64("r", 7092701, "room_id")
	flag.Parse()
	logrus.SetFormatter(&logrus.TextFormatter{})
	client, err := danmaku.Dial(danmaku.DouyuDanmakuServer, danmaku.WithRoom(*room))
	if err != nil {
		logrus.Errorf("dial danmu error: %v", err)
		os.Exit(1)
	}
	if err = client.Run(); err != nil {
		logrus.Errorf("exit with error: %v", err)
		os.Exit(2)
	}
	logrus.Info("exit successful")
}

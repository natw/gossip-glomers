package main

import (
	"context"
	"log"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {
	messageStore = make(map[int]bool)

	n := maelstrom.NewNode()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	writec := make(chan int)
	readc := make(chan chan []int)
	topWritec := make(chan []string)
	topReadc := make(chan chan []string)

	go messageMgr(ctx, readc, writec)
	go topologyMgr(ctx, topReadc, topWritec)

	n.Handle("broadcast", broadcastHandler(n, writec, topReadc))
	n.Handle("read", readHandler(n, readc))
	n.Handle("topology", topologyHandler(n, topWritec))

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}

type BroadcastReq struct {
	Type    string `json:"type"`
	Message int    `json:"message"`
}

type OKResp struct {
	Type string `json:"type"`
}

type ReadReq struct {
	Type string `json:"type"`
}

type ReadResp struct {
	Type     string `json:"type"`
	Messages []int  `json:"messages"`
}

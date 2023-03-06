package main

import (
	"log"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {
	n := maelstrom.NewNode()
	app := NewApp(n)

	n.Handle("broadcast", app.BroadcastHandler)
	n.Handle("broadcast_ok", func(msg maelstrom.Message) error {
		return nil
	})
	n.Handle("read", app.ReadHandler)
	n.Handle("topology", app.TopologyHandler)

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}

type BroadcastReq struct {
	Type    string   `json:"type"`
	SeenBy  []string `json:"seen_by"`
	Message int      `json:"message"`
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

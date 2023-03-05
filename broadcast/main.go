package main

import (
	"context"
	"encoding/json"
	"log"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

var messageStore map[int]bool

func main() {
	messageStore = make(map[int]bool)

	n := maelstrom.NewNode()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	writec := make(chan int)
	readc := make(chan chan []int)

	go messageMgr(ctx, readc, writec)

	n.Handle("broadcast", broadcastHandler(n, writec))
	n.Handle("read", readHandler(n, readc))
	n.Handle("topology", topologyHandler(n))

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}

func messageMgr(ctx context.Context, readc <-chan chan []int, writec <-chan int) {
	for {
		select {
		case <-ctx.Done():
			return
		case readrespc := <-readc:
			keys := make([]int, len(messageStore))
			i := 0
			for msg := range messageStore {
				keys[i] = msg
				i++
			}
			readrespc <- keys
		case msg := <-writec:
			messageStore[msg] = true
		}
	}
}

func broadcastHandler(n *maelstrom.Node, writec chan int) func(maelstrom.Message) error {
	return func(msg maelstrom.Message) error {
		var err error
		var req BroadcastReq
		err = json.Unmarshal(msg.Body, &req)
		if err != nil {
			return err
		}

		writec <- req.Message

		resp := OKResp{Type: "broadcast_ok"}
		return n.Reply(msg, resp)
	}
}

func readHandler(n *maelstrom.Node, readc chan chan []int) func(maelstrom.Message) error {
	return func(msg maelstrom.Message) error {
		var err error

		var req ReadReq
		err = json.Unmarshal(msg.Body, &req)
		if err != nil {
			return err
		}

		mc := make(chan []int, 1)

		readc <- mc

		msgs := <-mc

		resp := ReadResp{
			Type:     "read_ok",
			Messages: msgs,
		}
		return n.Reply(msg, resp)
	}
}

func topologyHandler(n *maelstrom.Node) func(maelstrom.Message) error {
	return func(msg maelstrom.Message) error {
		var err error

		var req TopologyReq
		err = json.Unmarshal(msg.Body, &req)
		if err != nil {
			return err
		}

		resp := OKResp{Type: "topology_ok"}
		return n.Reply(msg, resp)
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

type TopologyReq struct {
	Topology map[string][]string
	Type     string
}

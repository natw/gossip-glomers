package main

import (
	"context"
	"encoding/json"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

var messageStore map[int]bool

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

func broadcastHandler(n *maelstrom.Node, writec chan int, topReadc chan chan []string) func(maelstrom.Message) error {
	return func(msg maelstrom.Message) error {
		var err error
		var req BroadcastReq
		err = json.Unmarshal(msg.Body, &req)
		if err != nil {
			return err
		}

		writec <- req.Message

		go forwardBroadcast(n, req, topReadc)

		resp := OKResp{Type: "broadcast_ok"}
		return n.Reply(msg, resp)
	}
}

// send a Broadcast request on to every neighbor
func forwardBroadcast(n *maelstrom.Node, req BroadcastReq, topReadc chan chan []string) {
	nbrc := make(chan []string)
	go func(nbrc chan []string) {
		neighbors := <-nbrc
		for _, nbr := range neighbors {
			n.Send(nbr, req)
		}
	}(nbrc)
	topReadc <- nbrc
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

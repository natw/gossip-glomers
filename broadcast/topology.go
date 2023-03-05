package main

import (
	"context"
	"encoding/json"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

var neighbors []string

func topologyHandler(n *maelstrom.Node, topWritec chan<- []string) func(maelstrom.Message) error {
	return func(msg maelstrom.Message) error {
		var err error

		var req TopologyReq
		err = json.Unmarshal(msg.Body, &req)
		if err != nil {
			return err
		}

		topWritec <- req.Topology[n.ID()]

		resp := OKResp{Type: "topology_ok"}
		return n.Reply(msg, resp)
	}
}

type TopologyReq struct {
	Topology map[string][]string
	Type     string
}

func topologyMgr(ctx context.Context, topWritec chan chan []string, topReadc <-chan []string) {
	for {
		select {
		case <-ctx.Done():
			return
		case nbrs := <-topReadc:
			neighbors = nbrs
		case reqc := <-topWritec:
			reqc <- neighbors
		}
	}
}

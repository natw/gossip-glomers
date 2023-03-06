package main

import (
	"encoding/json"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func (a *App) TopologyHandler(msg maelstrom.Message) error {
	var err error

	var req TopologyReq
	err = json.Unmarshal(msg.Body, &req)
	if err != nil {
		return err
	}

	a.SetNeighbors(req.Topology[a.Node.ID()])

	resp := OKResp{Type: "topology_ok"}
	return a.Node.Reply(msg, resp)
}

func (a *App) SetNeighbors(nbrs []string) {
	a.neighbors = nbrs
}

func (a *App) GetNeighbors() []string {
	return a.neighbors
}

type TopologyReq struct {
	Topology map[string][]string
	Type     string
}

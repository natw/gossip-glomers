package main

import (
	"encoding/json"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func (a *App) BroadcastHandler(msg maelstrom.Message) error {
	var err error
	var req BroadcastReq
	err = json.Unmarshal(msg.Body, &req)
	if err != nil {
		return err
	}

	a.AddMessage(req.Message)

	go func(req BroadcastReq, nodeID string, src string) {
		req.SeenBy = append(req.SeenBy, nodeID)
		req.SeenBy = append(req.SeenBy, src)
		seen := make(map[string]bool)
		for _, n := range req.SeenBy {
			seen[n] = true
		}

		nbrs := a.GetNeighbors()
		for _, nbr := range nbrs {
			if !seen[nbr] {
				a.forwardBroadcast(nbr, req)
			}
		}
	}(req, a.Node.ID(), msg.Src)

	resp := OKResp{Type: "broadcast_ok"}
	return a.Node.Reply(msg, resp)
}

func (a *App) forwardBroadcast(dest string, req BroadcastReq) {
	a.Logf("node '%s' forwarding broadcast of '%v' to '%s'", a.Node.ID(), req, dest)
	_ = a.Node.Send(dest, req)
}

func (a *App) ReadHandler(msg maelstrom.Message) error {
	var err error

	var req ReadReq
	err = json.Unmarshal(msg.Body, &req)
	if err != nil {
		return err
	}

	resp := ReadResp{
		Type:     "read_ok",
		Messages: a.GetMessages(),
	}
	return a.Node.Reply(msg, resp)
}

func (a *App) GetMessages() []int {
	a.neighborLock.Lock()
	defer a.neighborLock.Unlock()

	keys := make([]int, len(a.messages))
	i := 0
	for msg := range a.messages {
		keys[i] = msg
		i++
	}
	return keys
}

func (a *App) AddMessage(msg int) {
	a.msgLock.Lock()
	defer a.msgLock.Unlock()
	a.messages[msg] = true
}

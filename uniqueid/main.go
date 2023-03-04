package main

import (
	"encoding/json"
	"fmt"
	"log"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {
	n := maelstrom.NewNode()

	n.Handle("generate", func(msg maelstrom.Message) error {
		var err error
		var req GenerateReq

		err = json.Unmarshal(msg.Body, &req)
		if err != nil {
			return err
		}

		resp := &GenerateOK{}
		resp.Type = "generate_ok"
		resp.Id = generateUniqueID()
		return n.Reply(msg, resp)
	})

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}

var x int

func generateUniqueID() string {
	x = x + 1
	return fmt.Sprint(x)
}

type GenerateReq struct {
	MsgId any    `json:"msg_id"`
	Type  string `json:"type"`
}

type GenerateOK struct {
	Type string `json:"type"`
	Id   string `json:"id"`
}

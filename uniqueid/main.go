package main

import (
	"encoding/json"
	"log"

	"github.com/google/uuid"

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

func generateUniqueID() string {
	return uuid.New().String()
}

type GenerateReq struct {
	MsgId any    `json:"msg_id"`
	Type  string `json:"type"`
}

type GenerateOK struct {
	Type string `json:"type"`
	Id   string `json:"id"`
}

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type BroadcastNode struct {
	node         *maelstrom.Node
	messagesSeen map[string]bool
	topology     map[string]([]string)
}

type Topology struct {
	Message_type string              `json:"type"`
	Topo         map[string][]string `json:"topology"`
}

type Broadcast struct {
	Message_type string `json:"type"`
	Message      int    `json:"message"`
}

type Read struct {
	Message_type string `json:"type"`
	Msg_id       int    `json:"msg_id"`
}

func main() {

	var bcn BroadcastNode
	bcn.node = maelstrom.NewNode()
	bcn.messagesSeen = make(map[string]bool)

	bcn.node.Handle("broadcast", func(msg maelstrom.Message) error {
		var body map[string]string
		var broadcast Broadcast
		if err := json.Unmarshal(msg.Body, &broadcast); err != nil {
			return err
		}
		fmt.Fprintf(os.Stderr, "here1\n")
		if _, seen := bcn.messagesSeen[strconv.Itoa(broadcast.Message)]; !seen {
			bcn.messagesSeen[strconv.Itoa(broadcast.Message)] = true
			for _, dest := range bcn.topology[bcn.node.ID()] {
				bcn.node.Send(dest, body)
			}
		}

		body = make(map[string]string)
		body["type"] = "broadcast_ok"

		return bcn.node.Reply(msg, body)
	})

	bcn.node.Handle("topology", func(msg maelstrom.Message) error {

		var body Topology
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}
		bcn.topology = body.Topo

		var rep map[string]string
		rep = make(map[string]string)
		rep["type"] = "topology_ok"
		return bcn.node.Reply(msg, rep)
	})

	bcn.node.Handle("read", func(msg maelstrom.Message) error {
		var body map[string]any
		body = make(map[string]any)
		var r Read
		if err := json.Unmarshal(msg.Body, &r); err != nil {
			return err
		}

		body["type"] = "read_ok"
		x := make([]int, 0)
		for key := range bcn.messagesSeen {
			k, _ := strconv.Atoi(key)
			x = append(x, k)
		}
		body["messages"] = x

		return bcn.node.Reply(msg, body)

	})

	if err := bcn.node.Run(); err != nil {
		log.Fatal((err))
	}
}

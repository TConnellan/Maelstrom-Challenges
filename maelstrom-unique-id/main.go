package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {
	// have nodes send request to another service
	// to retrieve pages of new unique ids every time they are running low

	node := maelstrom.NewNode()

	//max := -1
	//current_val := 0
	unique_ids := make(map[string][]int)

	node.Handle("generate", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		if len(unique_ids["ids"]) == 0 {
			if resp, err := http.Get("http://localhost:8080"); err != nil {
				log.Fatal(err)
			} else {
				if values, read_err := io.ReadAll(resp.Body); read_err != nil {
					log.Fatal(err)
				} else {
					json.Unmarshal(values, &unique_ids)
				}
				resp.Body.Close()
			}

		}

		body["type"] = "generate_ok"

		body["id"] = unique_ids["ids"][0]
		unique_ids["ids"] = unique_ids["ids"][1:]

		return node.Reply(msg, body)
	})

	if err := node.Run(); err != nil {
		log.Fatal((err))
	}
}

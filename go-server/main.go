package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	//"encoding/json"
)

func main() {

	current_max := 0
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var body map[string]any
		body = make(map[string]any)
		//body["ids"] = [100]int{current_max +1,...,current_max + 100}
		//i := current_max + 1

		body["ids"] = get_ids(&current_max)
		//body["ids"] = []int{1, 2, 3, 4, 5}
		data, err := json.Marshal(body)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Fprintf(w, fmt.Sprintf("%s", data))
		}

		current_max = current_max + 100
	})
	fmt.Println("Running go-server on port 8080")
	go doFancyWords()
	go log.Fatal(http.ListenAndServe(":8080", nil))
}

func get_ids(current_max *int) []int {
	x := make([]int, 100)

	for pos := range x {
		(*current_max)++
		x[pos] = *current_max
	}

	return x
}

func doFancyWords() {

	for i := 0; i < 1; {
		for _, char := range "running..." {
			time.Sleep(500 * time.Millisecond)
			fmt.Printf("%c", char)

		}
		time.Sleep(500 * time.Millisecond)
		// clear the line
		fmt.Print("\033[2K")
		// return to start of line
		fmt.Print("\r")
	}

}

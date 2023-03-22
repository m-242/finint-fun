package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	invest := Investigate("A", []string{"BSActor"}, 24, "BS", "")

	j, err := json.Marshal(invest)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(j))
}

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
	invest := Investigate(
		"EQBq_FU3HKy1NykhavHmxh05IBhM6X1u2Y489AcvaXdg29ue",
		[]string{"Killnet"},
		6,
		"TON",
		"",
	)

	//invest := Investigate(
	//	"A",
	//	[]string{"Killnet"},
	//	10,
	//	"BS",
	//	"",
	//)

	j, err := json.Marshal(invest)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(j))

	invest.ToSVG("/tmp/invest_test2.svg")
}

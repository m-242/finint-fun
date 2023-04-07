package main

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	//invest := Investigate(
	//	"EQBq_FU3HKy1NykhavHmxh05IBhM6X1u2Y489AcvaXdg29ue",
	//	[]string{"Killnet"},
	//	3,
	//	"TON",
	//	"",
	//)
	//invest := Investigate(
	//	"A",
	//	[]string{"Killnet"},
	//	10,
	//	"BS",
	//	"",
	//)

	invest := Investigate(
		"bc1qt8ae83jf98rdffw2vmhcgg9l8w3lqy844vhdnt",
		[]string{"Killnet"},
		3,
		"BTC",
		"",
	)

	// j, err := json.Marshal(invest)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println(string(j))

	//invest, _ := fromJSON("/tmp/invest.json")
	//invest.toJSON("/tmp/invest.json")
	invest.toCytoscapeJS("/tmp/invest.html")
	//invest.ToSVG("/tmp/invest_test2.svg")
}

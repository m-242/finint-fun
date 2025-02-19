package main

import (
	"fmt"
	// "log"
	"io/ioutil"
	"strings"

	_ "embed"
)

//go:embed cytoscape/index.tmpl
var cytoscape_output_template string

// yup, fucking disgusting
func (invest *Investigation) toCytoscapeJS(path string) error {
	minBalance, maxBalance := invest.getNodeMinMax()
	minNode, maxNode := 30.0, 300.0

	cytoscape_tmpl := strings.Replace(
		cytoscape_output_template,
		"INVEST_CURRENCY",
		invest.Currency,
		1,
	)

	array := "["

	for _, node := range invest.InvolvedAddresses {
		if len(node.Tags) == 0 {
			array += fmt.Sprintf(
				"{ data: { id: '%s', weight: %d, balance: %d } }",
				node.Identifier,
				int(node.Balance/(maxBalance-minBalance)*(maxNode-minNode)+minNode),
				int(node.Balance),
			)
		} else {
			array += fmt.Sprintf(
				"{ data: { id: '%s', weight: %d, balance: %d, style: { 'background-color': 'red'} } }",
				node.Identifier,
				int(node.Balance/(maxBalance-minBalance)*(maxNode-minNode)+minNode),
				int(node.Balance),
			)

		}

		array += ",\n"
	}

	rMax := len(invest.Transactions)
	for r, trans := range invest.Transactions {
		array += fmt.Sprintf(
			"{data: {source: '%s', target: '%s', shape:'arrow', label: '%s', value: %f, currency: '%s' }}",
			trans.From,
			trans.To,
			fmt.Sprintf("%.2f %s", trans.Value, invest.Currency),
			trans.Value,
			invest.Currency,
		)

		if r != rMax {
			array += ",\n"
		}
	}

	array += "]"

	// Let's build the json array
	cytoscape_tmpl = strings.Replace(cytoscape_tmpl, "INVEST_DATA", array, 1)

	//fmt.Println(cytoscape_tmpl)

	return ioutil.WriteFile(path, []byte(cytoscape_tmpl), 0644)
}

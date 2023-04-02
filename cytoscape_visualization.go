package main

import (
	"fmt"
	// "log"
	"io/ioutil"
	"strings"
)

import _ "embed"

//go:embed cytoscape/cytoscape.js
var cytoscape_source string

//go:embed cytoscape/index.tmpl
var cytoscape_output_template string

// yup, fucking disgusting
func (invest *Investigation) toCytoscapeJS(path string) error {
	minBalance, maxBalance := invest.getNodeMinMax()
	minNode, maxNode := 1.0, 10.0
	//
	// minValue, maxValue := invest.getVerticeMinMax()
	// minVertice, maxVertice := 3.0, 25.0

	cytoscape_tmpl := strings.Replace(
		cytoscape_output_template,
		"CYTOSCAPE_SOURCE",
		cytoscape_source,
		1,
	)

	cytoscape_tmpl = strings.Replace(
		cytoscape_tmpl,
		"INVEST_CURRENCY",
		invest.Currency,
		1,
	)

	array := "["

	for _, node := range invest.InvolvedAddresses {
		array += fmt.Sprintf(
			"{ data: { id: '%s', width: '%d' } }",
			node.Identifier,
			int(node.Balance/(maxBalance-minBalance)*(maxNode-minNode)+minNode),
		)

		array += ","
	}

	rMax := len(invest.Transactions)
	for r, trans := range invest.Transactions {
		array += fmt.Sprintf("{data: {source: '%s', target: '%s' }}", trans.From, trans.To)

		if r != rMax {
			array += ","
		}
	}

	array += "]"

	// Let's build the json array
	cytoscape_tmpl = strings.Replace(cytoscape_tmpl, "INVEST_DATA", array, 1)

	//fmt.Println(cytoscape_tmpl)

	return ioutil.WriteFile(path, []byte(cytoscape_tmpl), 0644)
}

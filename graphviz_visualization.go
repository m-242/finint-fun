package main

import (
	"fmt"
	"log"
	"math"

	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
)

// Node is address. Size of Node iś depending on balance
// Vertex is transaction. Thickness of vertex depends on value
func (invest *Investigation) ToSVG(path string) error {
	minBalance, maxBalance := invest.getNodeMinMax()
	minNode, maxNode := 1.0, 10.0

	minValue, maxValue := invest.getVerticeMinMax()
	minVertice, maxVertice := 3.0, 25.0

	g := graphviz.New()
	graph, err := g.Graph()
	if err != nil {
		log.Fatal(err)
	}

	nodes := make(map[string]*cgraph.Node)

	for _, node := range invest.InvolvedAddresses {
		x, _ := graph.CreateNode(node.Identifier)
		nodes[node.Identifier] = x

		// Pretty stuff
		x.SetShape("circle")
		x.SetWidth(node.Balance/(maxBalance-minBalance)*(maxNode-minNode) + minNode)

		// TODO handle tags properly
		if len(node.Tags) > 0 {
			x.SetStyle(cgraph.FilledNodeStyle)
			x.SetFillColor("red")
			x.SetColor("red")
		}
	}

	for _, edge := range invest.Transactions {

		x, ok := nodes[edge.From]
		if !ok {
			x, _ := graph.CreateNode(edge.From)
			x.SetShape("circle")
		}

		y, ok2 := nodes[edge.To]
		if !ok2 {
			y, _ := graph.CreateNode(edge.To)
			y.SetShape("circle")
		}

		v, err := graph.CreateEdge(
			fmt.Sprintf("%f", edge.Value),
			x,
			y,
		)
		if err != nil {
			invest.logger.Printf("Coul'nt print edge")
		} else {
			v.SetPenWidth(edge.Value/(maxValue-minValue)*(maxVertice-minVertice) + minVertice)
			v.SetLabel(fmt.Sprintf("%f %s", edge.Value, invest.Currency))
		}
	}

	if err := g.RenderFilename(graph, graphviz.SVG, path); err != nil {
		log.Fatal(err)
	}

	return nil
}

func (invest *Investigation) getVerticeMinMax() (float64, float64) {
	max := math.SmallestNonzeroFloat64
	min := math.MaxFloat64

	for _, t := range invest.Transactions {
		if t.Value > max {
			max = t.Value
		}
		if t.Value < min {
			min = t.Value
		}
	}

	return min, max
}

func (invest *Investigation) getNodeMinMax() (float64, float64) {
	max := math.SmallestNonzeroFloat64
	min := math.MaxFloat64

	for _, t := range invest.InvolvedAddresses {
		if t.Balance > max {
			max = t.Balance
		}
		if t.Balance < min {
			min = t.Balance
		}
	}

	return min, max
}

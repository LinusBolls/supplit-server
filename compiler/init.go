package compiler

import (
	"errors"
	"fmt"

	"../types"
)

type Primitive = types.Primitive
type Node = types.Node
type BodyNode = types.BodyNode
type InNode = types.InNode
type OutNode = types.OutNode

type Point = types.Point
type Schema = types.NodeMapSchema

type NodeMapParser struct {
	NodeTypes map[string]Node
	Schema    Schema
	Computed  []Primitive
	Nodes     []string
}

func (p NodeMapParser) MakeNodeStruct(nodeId Point) Node {
	nodeName := p.Nodes[nodeId.Node]
	bodyNode := p.NodeTypes[nodeName]

	isInNode := nodeId.Node < len(p.Schema.In.Columns)
	isOutNode := nodeId.Node >= len(p.Schema.In.Columns)+len(p.Schema.Nodes)

	if isInNode {
		fmt.Println(nodeName + " in")
		return InNode{Type: "In"}
	}
	if isOutNode {
		fmt.Println(nodeName + " out")
		return OutNode{Type: "out", In: []string{"sache"}}
	}
	if bodyNode != nil {

		fmt.Println(nodeName + " body")
		return bodyNode
	}

	errors.New("Failed to find node identity for nodeId")

	return nil
}
func (p NodeMapParser) ResolveNodeInput(nodeId Point) Primitive {

	var noodle [2][2]int
	isAssigned := false

	for _, n := range p.Schema.Noodles {
		if n[1][0] == nodeId.Node && n[1][1] == nodeId.Port {
			noodle = n
			isAssigned = true
		}
	}
	if !isAssigned {
		errors.New("Failed to find noodle connection for nodeId")
	}
	sache := Point{noodle[0][0], noodle[0][1]}
	return p.ResolveNode(sache)
}
func (p NodeMapParser) ResolveNode(nodeId Point) []Primitive {

	existingValue := p.Computed[nodeId.Node]

	if existingValue != nil {
		return []Primitive{existingValue}
	}

	nodeStruct := p.MakeNodeStruct(nodeId)

	// passed node should never be an in node
	if _, ok := nodeStruct.(InNode); ok {
		errors.New("Input node value was not found in input parameter of calculateOutNodes")
	}

	var inValues []Primitive

	// cast to OutNode only to access In property

	_, isBodyNode := nodeStruct.(BodyNode)

	if isBodyNode {
		for idx, _ := range nodeStruct.(BodyNode).In {

			inValues = append(inValues, p.ResolveNodeInput(Point{nodeId.Node, idx}))

			fmt.Println(inValues)
		}
	} else {
		for idx, _ := range nodeStruct.(OutNode).In {
			inValues = append(inValues, p.ResolveNodeInput(Point{nodeId.Node, idx}))

			fmt.Println(inValues)
		}
	}

	if isBodyNode {
		result := nodeStruct.(BodyNode).Calc(inValues)
		p.Computed[nodeId.Node] = result[0]

		return result
	} else {
		result := inValues
		p.Computed[nodeId.Node] = result

		return result
	}

}
func ParseSache(nodeTypes map[string]Node, schema Schema, input []Primitive) []Primitive {
	notOutNodes := append(schema.In.Columns, schema.Nodes...)
	notInNodes := append(schema.Nodes, schema.Out.Columns...)
	nodes := append(notOutNodes, schema.Out.Columns...)

	p := new(NodeMapParser)
	p.NodeTypes = nodeTypes
	p.Schema = schema
	empty := make([]Primitive, len(notInNodes))
	p.Computed = append(input, empty...)

	p.Nodes = nodes

	for idx, _ := range schema.Out.Columns {
		nodeId := Point{idx + len(notOutNodes), 0}

		p.ResolveNode(nodeId)
	}

	computedOutNodes := p.Computed[len(notOutNodes):]

	return computedOutNodes
}

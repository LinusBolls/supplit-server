package types

import (
	"fmt"
	"strconv"
)

type Primitive interface{}

type SchemaIn struct {
	Columns []string `json:"columns"`
}
type SchemaOut struct {
	Columns []string `json:"columns"`
}
type NodeMapSchema struct {
	In      SchemaIn    `json:"in"`
	Out     SchemaOut   `json:"out"`
	Nodes   []string    `json:"nodes"`
	Noodles [][2][2]int `json:"noodles"`
}

type Point struct {
	Node int
	Port int
}

type Node interface{}

type BodyNode struct {
	Type string

	In   []string
	Out  []string
	Calc func([]Primitive) []Primitive
}
type InNode struct {
	Type string
}
type OutNode struct {
	Type string
	In   []string
}

var DefaultNodeTypes = map[string]Node{
	"toPercent": BodyNode{
		"body",
		[]string{""},
		[]string{""},
		func(args []Primitive) []Primitive {
			return args
		},
	},
	"multiply": BodyNode{
		"body",
		[]string{"one", "two"},
		[]string{"result"},
		func(args []Primitive) []Primitive {

			firstFloat, err := strconv.ParseFloat(args[0].(string), 10)
			secondFloat, err := strconv.ParseFloat(args[1].(string), 10)

			if err != nil {
				fmt.Println(err)
			}
			resultFloat := firstFloat * secondFloat

			fmt.Println("result in bodynode:")
			fmt.Println([]Primitive{Primitive(resultFloat)})
			fmt.Println([]Primitive{Primitive(resultFloat)}[0])

			return []Primitive{Primitive(resultFloat)}
		},
	},
}

type NodeMapErrorSeverity struct {
	Name string `json:"name"`
	Desc string `json:"desc"`
	Icon string `json:"icon"` // a font awesome class

	PrimColor string `json:"primColor"` // a css color
	SecColor  string `json:"secColor"`  //a css color
}
type NodeMapError struct {
	Severity NodeMapErrorSeverity `json:"severity"`
	Source   string               `json:"source"` // "Ok" | "Compiler" | "Schema" | "Input Csv" | "Output Csv";
	Desc     string               `json:"desc"`
}

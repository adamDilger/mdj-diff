package main

import "encoding/json"

type Project struct{ Node }

func (n *Node) UnmarshalJSON(data []byte) error {
	var hello struct {
		N_type        string `json:"_type"`
		Id            string `json:"_id"`
		Name          string
		OwnedElements []*json.RawMessage
	}

	if err := json.Unmarshal(data, &hello); err != nil {
		return err
	}

	*n = Node{N_type: hello.N_type, Id: hello.Id, Name: hello.Name}

	for _, i := range hello.OwnedElements {
		var a struct {
			N_type string `json:"_type"`
		}

		if err := json.Unmarshal(*i, &a); err != nil {
			panic(err)
		}

		var dst any
		switch a.N_type {
		case "ERDDataModel":
			dst = new(DataModel)
		case "ERDEntity":
			dst = new(Entity)
		case "ERDDiagram":
			dst = new(Diagram)
		default:
			println(a.N_type)
			dst = new(interface{})
		}

		if err := json.Unmarshal(*i, dst); err != nil {
			panic(err)
		}

		(*n).OwnedElements = append((*n).OwnedElements, dst)
	}

	return nil
}

type Node struct {
	N_type        string `json:"_type"`
	Id            string `json:"_id"`
	Name          string `json:"name"`
	OwnedElements []any
}

type Element struct {
	Node
	Parent        Ref    `json:"_parent"`
	Tags          []Tag  `json:"tags"`
	Documentation string `json:"documentation"`
}

type Ref struct {
	Ref string `json:"$ref"`
}

type Diagram struct {
	Element
	DefaultDiagram bool `json:"defaultDiagram"`
}

type DataModel struct{ Element }
type Entity struct {
	Element
	Columns []Column `json:"columns"`
}

type Tag struct {
	Element
	Kind  string `json:"kind"`
	Value string `json:"value"`
}

type Column struct {
	Element
	Type        string `json:"type"`
	ReferenceTo Ref    `json:"referenceTo"`
	PrimaryKey  bool   `json:"primaryKey"`
	ForeignKey  bool   `json:"foreignKey"`
	Length      string `json:"length"`
}

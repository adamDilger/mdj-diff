package types

import (
	"encoding/json"
	"io"
)

type Base struct {
	Id            string      `json:"_id"`
	Name          string      `json:"name"`
	Parent        Ref         `json:"_parent"`
	OwnedElements ElementList `json:"ownedElements"`
	Tags          []Tag       `json:"tags"`
	Documentation string      `json:"documentation"`
}

func (b *Base) GetId() string                 { return b.Id }
func (b *Base) GetName() string               { return b.Name }
func (b *Base) GetParent() Ref                { return b.Parent }
func (b *Base) GetOwnedElements() ElementList { return b.OwnedElements }
func (b *Base) GetTags() []Tag                { return b.Tags }
func (b *Base) GetDocumentation() string      { return b.Documentation }

type Node interface {
	GetNodeType() string
	GetId() string
	GetName() string
	GetOwnedElements() ElementList
	GetParent() Ref
	GetTags() []Tag
	GetDocumentation() string
}

type ElementList []Node

var _ Node = (*Project)(nil)

type Project struct {
	Base
	RefLookup map[string]Node
	diagrams  []Node
}

func NewProjectFromJson(in io.Reader) (*Project, error) {
	var project *Project
	if err := json.NewDecoder(in).Decode(&project); err != nil {
		return nil, err
	}

	project.RefLookup = make(map[string]Node)

	elements := project.GetElementsOfType(ENTITY)
	for _, e := range elements {
		if en, ok := e.(*Entity); ok {
			project.RefLookup[e.GetId()] = en

			for _, c := range en.Columns {
				project.RefLookup[c.GetId()] = &c
			}
		}
	}

	project.diagrams = project.GetElementsOfType(DIAGRAM)

	return project, nil
}

func (p *Project) GetNodeType() string { return "Project" }

func (p *Project) GetTableMap() map[string]Entity {
	out := make(map[string]Entity)
	entities := p.GetElementsOfType(ENTITY)

	for _, e := range entities {
		if en, ok := e.(*Entity); ok {
			out[e.GetId()] = *en
		} else {
			panic(e.GetId() + "is not an entity")
		}
	}

	return out
}

func (p *Project) GetElementsOfType(typ string) []Node {
	return findElements(p.OwnedElements, typ, &[]Node{})
}

func findElements(el ElementList, typ string, nodes *[]Node) []Node {
	if len(el) == 0 {
		return *nodes
	}

	for _, e := range el {
		if e.GetNodeType() == typ {
			*nodes = append(*nodes, e)
		}

		findElements(e.GetOwnedElements(), typ, nodes)
	}

	return *nodes
}

type Element struct {
	Base
	N_type string `json:"_type"`
}

func (e *Element) GetNodeType() string { return e.N_type }

func (e *ElementList) UnmarshalJSON(data []byte) error {
	var allElements []*json.RawMessage
	if err := json.Unmarshal(data, &allElements); err != nil {
		return err
	}

	*e = ElementList{}

	for _, i := range allElements {
		var typ struct {
			Type string `json:"_type"`
		}
		if err := json.Unmarshal(*i, &typ); err != nil {
			panic(err)
		}

		var dst Node
		switch typ.Type {
		case RELATIONSHIP:
			dst = new(Relationship)
		case ENTITY:
			dst = new(Entity)
		case DIAGRAM:
			dst = new(Diagram)
		default:
			dst = new(Element)
		}

		if err := json.Unmarshal(*i, dst); err != nil {
			panic(err)
		}

		*e = append(*e, dst)
	}

	return nil
}

package types

import (
	"encoding/json"
	"io"
)

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
	Id            string `json:"_id"`
	Name          string `json:"name"`
	OwnedElements ElementList
	Parent        Ref    `json:"_parent"`
	Tags          []Tag  `json:"tags"`
	Documentation string `json:"documentation"`

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

func (p *Project) GetNodeType() string           { return "Project" }
func (p *Project) GetId() string                 { return p.Id }
func (p *Project) GetName() string               { return p.Name }
func (p *Project) GetOwnedElements() ElementList { return p.OwnedElements }
func (p *Project) GetParent() Ref                { return p.Parent }
func (p *Project) GetTags() []Tag                { return p.Tags }
func (p *Project) GetDocumentation() string      { return p.Documentation }

func (p *Project) GetTableMap() map[string]Entity {
	out := make(map[string]Entity)
	entities := p.GetElementsOfType(ENTITY)

	for _, e := range entities {
		if en, ok := e.(*Entity); ok {
			out[e.GetId()] = *en
		} else {
			panic(e.GetName() + "is not an entity")
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
	N_type        string `json:"_type"`
	Id            string `json:"_id"`
	Name          string
	OwnedElements ElementList
	Parent        Ref `json:"_parent"`
	Tags          []Tag
	Documentation string
}

func (e *Element) GetNodeType() string           { return e.N_type }
func (e *Element) GetId() string                 { return e.Id }
func (e *Element) GetName() string               { return e.Name }
func (e *Element) GetOwnedElements() ElementList { return e.OwnedElements }
func (e *Element) GetParent() Ref                { return e.Parent }
func (e *Element) GetTags() []Tag                { return e.Tags }
func (e *Element) GetDocumentation() string      { return e.Documentation }

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

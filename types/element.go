package types

import (
	"encoding/json"
	"fmt"
)

var _ Node = (*Element)(nil)

const (
	DIAGRAM          = "ERDDiagram"
	DATA_MODEL       = "ERDDataModel"
	PROJECT          = "Project"
	ENTITY           = "ERDEntity"
	COLUMN           = "ERDColumn"
	RELATIONSHIP     = "ERDRelationship"
	RELATIONSHIP_END = "ERDRelationshipEnd"
)

var _ Node = (*Entity)(nil)

type Entity struct {
	Base
	Columns []Column
}

func (e *Entity) GetNodeType() string { return ENTITY }

func (e *Entity) GetColumnMap() map[string]Column {
	out := make(map[string]Column)

	for _, c := range e.Columns {
		out[c.GetId()] = c
	}

	return out
}
func (e *Entity) GetRelationshipMap() map[string]Relationship {
	out := make(map[string]Relationship)

	for _, o := range e.OwnedElements {
		if r, ok := o.(*Relationship); ok {
			out[o.GetId()] = *r
		}
	}

	return out
}
func (e *Entity) GetTagMap() map[string]Tag {
	out := make(map[string]Tag)

	for _, t := range e.Tags {
		out[t.Name] = t
	}

	return out
}

type Ref struct {
	Ref string `json:"$ref"`
}

var _ Node = (*Diagram)(nil)

type Diagram struct {
	Base
	DefaultDiagram bool `json:"defaultDiagram"`
}

func (d *Diagram) GetNodeType() string { return DIAGRAM }

var _ Node = (*Column)(nil)

type Column struct {
	Base

	Type        string
	ReferenceTo Ref
	PrimaryKey  bool
	ForeignKey  bool
	Nullable    bool
	Unique      bool
	Length      ColumnLength
}

func (c *Column) GetNodeType() string { return COLUMN }

type ColumnLength string

func (c *ColumnLength) UnmarshalJSON(b []byte) error {
	var n interface{}
	if err := json.Unmarshal(b, &n); err != nil {
		return err
	}

	switch v := n.(type) {
	case string:
		*c = ColumnLength(v)
	case float64:
		*c = ColumnLength(fmt.Sprintf("%d", int(v)))
	default:
		return fmt.Errorf("unexpected type %T for value %v", v, n)
	}

	return nil
}

var _ Node = (*Relationship)(nil)

type Relationship struct {
	Base
	End1 RelationshipEnd
	End2 RelationshipEnd
}

func (r *Relationship) GetNodeType() string { return RELATIONSHIP }

type RelationshipEnd struct {
	Id          string `json:"_id"`
	Parent      Ref    `json:"_parent"`
	Reference   Ref
	Cardinality string
}

func (r *RelationshipEnd) GetCardinality() string {
	if r.Cardinality == "" {
		return "1"
	}

	return r.Cardinality
}

type Tag struct {
	Id            string `json:"_id"`
	Name          string
	Parent        Ref `json:"_parent"`
	Documentation string

	Kind  string `json:"kind"`
	Value string `json:"value"`
}

func GetTagMap(tags []Tag) map[string]Tag {
	out := make(map[string]Tag)

	for _, t := range tags {
		out[t.Name] = t
	}

	return out
}

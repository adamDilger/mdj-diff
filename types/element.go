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
	Id            string `json:"_id"`
	Name          string
	OwnedElements ElementList
	Parent        Ref `json:"_parent"`
	Tags          []Tag
	Documentation string

	Columns []Column
}

func (e *Entity) GetNodeType() string           { return ENTITY }
func (e *Entity) GetId() string                 { return e.Id }
func (e *Entity) GetName() string               { return e.Name }
func (e *Entity) GetOwnedElements() ElementList { return e.OwnedElements }
func (e *Entity) GetParent() Ref                { return e.Parent }
func (e *Entity) GetTags() []Tag                { return e.Tags }
func (e *Entity) GetDocumentation() string      { return e.Documentation }
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
	Id            string `json:"_id"`
	Name          string
	OwnedElements ElementList
	Parent        Ref `json:"_parent"`
	Tags          []Tag
	Documentation string

	DefaultDiagram bool `json:"defaultDiagram"`
}

func (d *Diagram) GetNodeType() string           { return DIAGRAM }
func (d *Diagram) GetId() string                 { return d.Id }
func (d *Diagram) GetName() string               { return d.Name }
func (d *Diagram) GetOwnedElements() ElementList { return d.OwnedElements }
func (d *Diagram) GetParent() Ref                { return d.Parent }
func (d *Diagram) GetTags() []Tag                { return d.Tags }
func (d *Diagram) GetDocumentation() string      { return d.Documentation }

var _ Node = (*Column)(nil)

type Column struct {
	Id            string `json:"_id"`
	Name          string
	Parent        Ref `json:"_parent"`
	Tags          []Tag
	Documentation string

	Type        string
	ReferenceTo Ref
	PrimaryKey  bool
	ForeignKey  bool
	Nullable    bool
	Unique      bool
	Length      ColumnLength
}

func (c *Column) GetNodeType() string           { return COLUMN }
func (c *Column) GetId() string                 { return c.Id }
func (c *Column) GetName() string               { return c.Name }
func (c *Column) GetOwnedElements() ElementList { return nil }
func (c *Column) GetParent() Ref                { return c.Parent }
func (c *Column) GetTags() []Tag                { return c.Tags }
func (c *Column) GetDocumentation() string      { return c.Documentation }

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
	Id            string `json:"_id"`
	Parent        Ref    `json:"_parent"`
	Tags          []Tag
	Documentation string

	End1 RelationshipEnd
	End2 RelationshipEnd
}

func (r *Relationship) GetNodeType() string           { return RELATIONSHIP }
func (r *Relationship) GetId() string                 { return r.Id }
func (r *Relationship) GetName() string               { return r.Parent.Ref }
func (r *Relationship) GetOwnedElements() ElementList { return nil }
func (r *Relationship) GetParent() Ref                { return r.Parent }
func (r *Relationship) GetTags() []Tag                { return r.Tags }
func (r *Relationship) GetDocumentation() string      { return r.Documentation }

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

package diff

import (
	"mdj-diff/types"
	"testing"
)

func TestDiffColumnWithNoChanges(t *testing.T) {
	a := types.Column{Id: "1", Name: "HELLO", Type: "INTEGER", Documentation: "doc"}
	b := types.Column{Id: "1", Name: "HELLO", Type: "INTEGER", Documentation: "doc"}

	cc := diffColumn(a, b)
	if cc != nil {
		t.Fatal("changes should be nil!")
	}
}

func TestDiffColumnModifyAttributes(t *testing.T) {
	a := types.Column{Id: "1", Name: "HELLO", Type: "INTEGER", Documentation: "doc"}
	b := types.Column{Id: "1", Name: "OTHER_TABLE", Type: "HELLO", Documentation: "doc1"}

	cc := diffColumn(a, b)
	if len(cc.Changes) != 3 {
		t.Fatalf("wanted %d columns, got %d", 3, len(cc.Changes))
	}

	// name attribute
	change := findChangeForField(t, cc.Changes, "name")
	if change.Type != ChangeTypeModify {
		t.Fatalf("wanted type %s, got %s", ChangeTypeModify, change.Type)
	}

	if change.Value != a.Name {
		t.Fatalf("wanted name: %s got %s", a.Name, change.Value)
	}

	if change.Old != b.Name {
		t.Fatalf("wanted name: %s got %s", b.Name, change.Old)
	}

	// type attribute
	change = findChangeForField(t, cc.Changes, "type")
	if change.Type != ChangeTypeModify {
		t.Fatalf("wanted type %s, got %s", ChangeTypeModify, change.Type)
	}

	if change.Value != a.Type {
		t.Fatalf("wanted name: %s got %s", a.Type, change.Value)
	}

	if change.Old != b.Type {
		t.Fatalf("wanted name: %s got %s", b.Type, change.Old)
	}

	// documentation attribute
	change = findChangeForField(t, cc.Changes, "documentation")
	if change.Type != ChangeTypeModify {
		t.Fatalf("wanted type %s, got %s", ChangeTypeModify, change.Type)
	}

	if change.Value != a.Documentation {
		t.Fatalf("wanted name: %s got %s", a.Documentation, change.Value)
	}

	if change.Old != b.Documentation {
		t.Fatalf("wanted name: %s got %s", b.Documentation, change.Old)
	}
}

func TestDiffColumnAddTag(t *testing.T) {
	tag := types.Tag{Name: "tag1", Value: "value1"}
	a := types.Column{Id: "1", Tags: []types.Tag{tag}}
	b := types.Column{Id: "1"}

	cc := diffColumn(a, b)
	if len(cc.Changes) != 0 {
		t.Fatalf("wanted %d changes, got %d", 0, len(cc.Changes))
	}
	if len(cc.Tags) != 1 {
		t.Fatalf("wanted %d tags, got %d", 1, len(cc.Tags))
	}

	change := cc.Tags[0]
	if change.Type != ChangeTypeAdd {
		t.Fatalf("wanted type %s, got %s", ChangeTypeModify, change.Type)
	}

	if change.Name != tag.Name {
		t.Fatalf("wanted name: %s got %s", tag.Name, change.Name)
	}

	if len(change.Changes) != 1 {
		t.Fatalf("wanted %d changes, got %d", 1, len(change.Changes))
	}

	c := change.Changes[0]
	if c.Type != ChangeTypeAdd {
		t.Fatalf("wanted type %s, got %s", ChangeTypeAdd, change.Type)
	}

	if c.Name != "value" && c.Value != tag.Value {
		t.Fatalf("wanted name: %s got %s", tag.Name, c.Name)
	}
}

func findChangeForField(t *testing.T, changes []Change, name string) Change {
	for _, c := range changes {
		if c.Name == name {
			return c
		}
	}

	t.Fatalf("could not find change for field %s", name)

	return Change{}
}

package diff

import (
	"mdj-diff/types"
	"testing"
)

func TestDiffEntityAttributeChange(t *testing.T) {
	b := types.Entity{Name: "name", Documentation: "doc"}
	a := types.Entity{Name: b.Name + "_updated", Documentation: b.Documentation + "_updated"}

	tc := diffEntity(a, b)
	if len(tc.Changes) != 2 {
		t.Fatalf("wanted %d changes, got %d", 2, len(tc.Changes))
	}

	for _, c := range tc.Changes {
		if c.Type != ChangeTypeModify {
			t.Fatalf("wanted type %s, got %s", ChangeTypeModify, c.Type)
		}

		if c.Name == "name" {
			if c.Old != b.Name {
				t.Errorf("wanted %s got %s", b.Name, c.Old)
			}
			if c.Value != a.Name {
				t.Errorf("wanted %s got %s", a.Name, c.Value)
			}
		}
		if c.Name == "documentation" {
			if c.Old != b.Documentation {
				t.Errorf("wanted %s got %s", b.Documentation, c.Old)
			}
			if c.Value != a.Documentation {
				t.Errorf("wanted %s got %s", a.Documentation, c.Value)
			}
		}
	}
}

func TestDiffEntityAddColumn(t *testing.T) {
	existingCol := types.Column{Id: "1", Name: "HELLO"}
	newCol := types.Column{Id: "2", Name: "OTHER_TABLE", Type: "HELLO"}
	b := types.Entity{Columns: []types.Column{existingCol}}
	a := types.Entity{Columns: []types.Column{existingCol, newCol}}

	tc := diffEntity(a, b)
	if len(tc.Columns) != 1 {
		t.Fatalf("wanted %d columns, got %d", 2, len(tc.Columns))
	}

	col := tc.Columns[0]
	if col.Type != ChangeTypeAdd {
		t.Fatalf("wanted type %s, got %s", ChangeTypeAdd, col.Type)
	}

	if col.Name != newCol.Name {
		t.Fatalf("wanted name: %s got %s", newCol.Name, col.Name)
	}

	if len(col.Changes) != 1 {
		t.Fatalf("wanted changes %d got %d", 1, len(col.Changes))
	}

	if col.Changes[0].Value != newCol.Type &&
		col.Changes[0].Name != "name" {
		t.Errorf("wanted %s got %s", newCol.Type, col.Changes[0].Type)
	}
}

func TestDiffEntityRemoveColumn(t *testing.T) {
	existingCol := types.Column{Id: "1", Name: "HELLO"}
	removeCol := types.Column{Id: "2", Name: "OTHER_TABLE", Type: "HELLO"}
	b := types.Entity{Columns: []types.Column{existingCol, removeCol}}
	a := types.Entity{Columns: []types.Column{existingCol}}

	tc := diffEntity(a, b)
	if len(tc.Columns) != 1 {
		t.Fatalf("wanted %d columns, got %d", 2, len(tc.Columns))
	}

	col := tc.Columns[0]
	if col.Type != ChangeTypeRemove {
		t.Fatalf("wanted type %s, got %s", ChangeTypeRemove, col.Type)
	}

	if col.Name != removeCol.Name {
		t.Fatalf("wanted name: %s got %s", removeCol.Name, col.Name)
	}

	if len(col.Changes) != 1 {
		t.Fatalf("wanted changes %d got %d", 1, len(col.Changes))
	}

	if col.Changes[0].Value != removeCol.Type &&
		col.Changes[0].Name != "name" {
		t.Errorf("wanted %s got %s", removeCol.Type, col.Changes[0].Type)
	}
}

func TestDiffEntityAddTag(t *testing.T) {
	existingCol := types.Column{Id: "1", Name: "HELLO"}
	colWithTag := types.Column{Id: "1", Name: "HELLO", Tags: []types.Tag{{Id: "123", Name: "Unique"}}}
	b := types.Entity{Columns: []types.Column{existingCol}}
	a := types.Entity{Columns: []types.Column{colWithTag}}

	tc := diffEntity(a, b)
	if len(tc.Columns) != 1 {
		t.Fatalf("wanted %d columns, got %d", 2, len(tc.Columns))
	}

	col := tc.Columns[0]
	if col.Type != ChangeTypeModify {
		t.Fatalf("wanted type %s, got %s", ChangeTypeModify, col.Type)
	}

	if len(col.Tags) != 1 {
		t.Fatalf("wanted changes %d got %d", 1, len(col.Tags))
	}

	tag := col.Tags[0]
	if tag.Name != colWithTag.Tags[0].Name {
		t.Fatalf("wanted name: %s got %s", col.Tags[0].Name, colWithTag.Tags[0].Name)
	}
}

func TestDiffEntityRemoveTag(t *testing.T) {
	existingCol := types.Column{Id: "1", Name: "HELLO"}
	colWithTag := types.Column{Id: "1", Name: "HELLO", Tags: []types.Tag{{Id: "123", Name: "Unique", Value: "true"}}}
	b := types.Entity{Columns: []types.Column{colWithTag}}
	a := types.Entity{Columns: []types.Column{existingCol}}

	tc := diffEntity(a, b)
	if len(tc.Columns) != 1 {
		t.Fatalf("wanted %d columns, got %d", 2, len(tc.Columns))
	}

	col := tc.Columns[0]
	if col.Type != ChangeTypeModify {
		t.Fatalf("wanted type %s, got %s", ChangeTypeModify, col.Type)
	}

	if len(col.Tags) != 1 {
		t.Fatalf("wanted changes %d got %d", 1, len(col.Tags))
	}

	tag := col.Tags[0]
	if tag.Name != colWithTag.Tags[0].Name {
		t.Fatalf("wanted name: %s got %s", col.Tags[0].Name, colWithTag.Tags[0].Name)
	}

	if tag.Changes[0].Name != "Unique" &&
		tag.Changes[0].Value != "true" {
		t.Fatalf("invalid tag changeset")
	}
}

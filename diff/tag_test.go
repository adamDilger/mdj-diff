package diff

import (
	"mdj-diff/types"
	"testing"
)

func TestDiffTagWithNoChanges(t *testing.T) {
	a := types.Tag{Name: "HELLO", Value: "INTEGER", Documentation: "doc"}
	b := types.Tag{Name: "HELLO", Value: "INTEGER", Documentation: "doc"}

	cc := diffTag(a, b)
	if cc != nil {
		t.Fatal("changes should be nil!")
	}
}

func TestDiffTagModifyAttributes(t *testing.T) {
	a := types.Tag{Name: "Default", Value: "'Y'", Documentation: "doc1", Kind: "two"}
	b := types.Tag{Name: "Default", Value: "'N'", Documentation: "doc", Kind: "one"}

	ct := diffTag(a, b)
	if len(ct.Changes) != 3 {
		t.Fatalf("wanted %d columns, got %d", 3, len(ct.Changes))
	}

	// value attribute
	change := findChangeForField(t, ct.Changes, "value")
	if change.Type != ChangeTypeModify {
		t.Fatalf("wanted type %s, got %s", ChangeTypeModify, change.Type)
	}

	if change.Value != a.Value {
		t.Fatalf("wanted name: %s got %s", a.Value, change.Value)
	}

	if change.Old != b.Value {
		t.Fatalf("wanted name: %s got %s", b.Value, change.Old)
	}

	// documentation attribute
	change = findChangeForField(t, ct.Changes, "documentation")
	if change.Type != ChangeTypeModify {
		t.Fatalf("wanted type %s, got %s", ChangeTypeModify, change.Type)
	}

	if change.Value != a.Documentation {
		t.Fatalf("wanted name: %s got %s", a.Documentation, change.Value)
	}

	if change.Old != b.Documentation {
		t.Fatalf("wanted name: %s got %s", b.Documentation, change.Old)
	}

	// kind attribute
	change = findChangeForField(t, ct.Changes, "kind")
	if change.Type != ChangeTypeModify {
		t.Fatalf("wanted type %s, got %s", ChangeTypeModify, change.Type)
	}

	if change.Value != a.Kind {
		t.Fatalf("wanted name: %s got %s", a.Kind, change.Value)
	}

	if change.Old != b.Kind {
		t.Fatalf("wanted name: %s got %s", b.Kind, change.Old)
	}
}

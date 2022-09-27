package diff

import (
	"fmt"
	"mdj-diff/types"
)

func diffColumns(a types.Entity, b types.Entity) []ColumnChange {
	var changes []ColumnChange

	// get column map
	aColumns := a.GetColumnMap()
	bColumns := b.GetColumnMap()

	existingMap := make(map[string]bool)

	for id, aCol := range aColumns {
		if bCol, ok := bColumns[id]; ok {
			existingMap[id] = true
			if cc := diffColumn(aCol, bCol); cc != nil {
				changes = append(changes, *cc)
			}
		} else {
			// new A column
			cc := wholeColumnChange(aCol, ChangeTypeAdd)
			changes = append(changes, cc)
		}
	}

	for id, bCol := range bColumns {
		if _, ok := existingMap[id]; ok {
			continue // already been diffed
		}

		cc := wholeColumnChange(bCol, ChangeTypeRemove)
		changes = append(changes, cc)
	}

	return changes
}

func diffColumn(a types.Column, b types.Column) *ColumnChange {
	cc := &ColumnChange{Id: a.Id, Name: a.GetName(), Type: ChangeTypeModify}

	fields := []struct {
		name string
		a, b interface{}
	}{
		{"name", a.GetName(), b.GetName()},
		{"documentation", a.GetDocumentation(), b.GetDocumentation()},
		{"type", a.Type, b.Type},
		{"primaryKey", a.PrimaryKey, b.PrimaryKey},
		{"foreignKey", a.ForeignKey, b.ForeignKey},
		{"nullable", a.Nullable, b.Nullable},
		{"unique", a.Unique, b.Unique},
		{"length", a.Length, b.Length},
	}

	for _, f := range fields {
		if f.a != f.b {
			cc.Changes = append(cc.Changes, Change{
				Name:  f.name,
				Type:  ChangeTypeModify,
				Value: fmt.Sprintf("%v", f.a),
				Old:   fmt.Sprintf("%v", f.b),
			})
		}
	}

	// for tags
	cc.Tags = diffTags(a.GetTags(), b.GetTags())

	if len(cc.Changes)+len(cc.Tags) == 0 {
		return nil
	}

	return cc
}

func wholeColumnChange(c types.Column, changeType ChangeType) ColumnChange {
	cc := ColumnChange{Id: c.Id, Name: c.GetName(), Type: changeType}

	cc.Changes = append(cc.Changes, Change{Name: "type", Type: changeType, Value: c.Type})

	// optional column fields
	if c.GetDocumentation() != "" {
		cc.Changes = append(cc.Changes, Change{Name: "documentation", Type: changeType, Value: c.GetDocumentation()})
	}

	if c.Length != "" {
		cc.Changes = append(cc.Changes, Change{Name: "length", Type: changeType, Value: string(c.Length)})
	}

	if c.PrimaryKey {
		cc.Changes = append(cc.Changes, Change{Name: "primaryKey", Type: changeType, Value: true})
	}
	if c.ForeignKey {
		cc.Changes = append(cc.Changes, Change{Name: "foreignKey", Type: changeType, Value: true})
	}

	return cc
}

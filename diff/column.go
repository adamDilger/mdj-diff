package diff

import (
	"mdj-diff/types"
)

func diffColumns(a types.Entity, b types.Entity) []BaseChange {
	var changes []BaseChange

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

func diffColumn(a types.Column, b types.Column) *BaseChange {
	cc := &BaseChange{Id: a.Id, Name: a.Name, Type: ChangeTypeModify}

	fields := []fieldDiff{
		{"name", a.Name, b.Name},
		{"documentation", a.Documentation, b.Documentation},
		{"type", a.Type, b.Type},
		{"primaryKey", a.PrimaryKey, b.PrimaryKey},
		{"foreignKey", a.ForeignKey, b.ForeignKey},
		{"nullable", a.Nullable, b.Nullable},
		{"unique", a.Unique, b.Unique},
		{"length", a.Length, b.Length},
	}

	cc.Changes = diffFields(fields)
	cc.Tags = diffTags(a.Tags, b.Tags)

	if len(cc.Changes)+len(cc.Tags) == 0 {
		return nil
	}

	return cc
}

func wholeColumnChange(c types.Column, changeType ChangeType) BaseChange {
	cc := BaseChange{Id: c.Id, Name: c.Name, Type: changeType}

	cc.Changes = append(cc.Changes, Change{Name: "type", Type: changeType, Value: c.Type})

	fields := []fieldDiff{
		{"name", c.Name, ""},
		{"documentation", c.Documentation, ""},
		{"primaryKey", c.PrimaryKey, false},
		{"foreignKey", c.ForeignKey, false},
		{"nullable", c.Nullable, false},
		{"unique", c.Unique, false},
		{"length", string(c.Length), ""},
	}

	cc.Changes = diffFieldsWithType(fields, changeType)

	return cc
}

package diff

import (
	"mdj-diff/types"
	"sort"
)

func DiffTables(A, B map[string]types.Entity) []TableChange {
	existingMap := make(map[string]bool)
	tableChanges := make([]TableChange, 0)

	for id, e := range A {
		bEntity, ok := B[id]
		if ok {
			existingMap[id] = true
			if tc := diffEntity(e, bEntity); tc != nil {
				tableChanges = append(tableChanges, *tc)
			}
		} else {
			// new A table
			tc := wholeTableChange(e, ChangeTypeAdd)
			tableChanges = append(tableChanges, tc)
		}
	}

	for id, e := range B {
		if _, ok := existingMap[id]; ok {
			continue // already been diffed
		}

		// new table in master, so mark as "removed" in the diff
		tc := wholeTableChange(e, ChangeTypeRemove)
		tableChanges = append(tableChanges, tc)
	}

	sort.Slice(tableChanges, func(i, j int) bool {
		return tableChanges[i].Name < tableChanges[j].Name
	})

	return tableChanges
}

func diffEntity(a types.Entity, b types.Entity) *TableChange {
	tc := &TableChange{Id: a.Id, Name: a.GetName(), Type: ChangeTypeModify}

	if a.GetName() != b.GetName() {
		c := Change{Name: "name", Type: ChangeTypeModify, Value: a.GetName(), Old: b.GetName()}
		tc.Changes = append(tc.Changes, c)
	}

	if a.GetDocumentation() != b.GetDocumentation() {
		c := Change{Name: "documentation", Type: ChangeTypeModify, Value: a.GetDocumentation(), Old: b.GetDocumentation()}
		tc.Changes = append(tc.Changes, c)
	}

	tc.Columns = diffColumns(a, b)
	tc.Relationships = diffRelationships(a, b)
	tc.Tags = diffTags(a.GetTags(), b.GetTags())

	if len(tc.Changes)+
		len(tc.Columns)+
		len(tc.Relationships)+
		len(tc.Tags) == 0 {
		return nil
	}

	return tc
}

func wholeTableChange(e types.Entity, changeType ChangeType) TableChange {
	tc := TableChange{Id: e.Id, Name: e.GetName(), Type: changeType}

	// optional table fields
	if e.GetDocumentation() != "" {
		tc.Changes = append(tc.Changes, Change{Name: "documentation", Type: changeType, Value: e.GetDocumentation()})
	}

	for _, col := range e.Columns {
		cc := wholeColumnChange(col, changeType)
		tc.Columns = append(tc.Columns, cc)
	}

	for _, tag := range e.Tags {
		cc := wholeTagChange(tag, changeType)
		tc.Tags = append(tc.Tags, cc)
	}

	return tc
}

package diff

import (
	"mdj-diff/types"
)

func diffTags(a []types.Tag, b []types.Tag) []BaseChange {
	var changes []BaseChange

	// get column map
	aTags := types.GetTagMap(a)
	bTags := types.GetTagMap(b)

	existingMap := make(map[string]bool)

	for id, aTag := range aTags {
		if bTag, ok := bTags[id]; ok {
			existingMap[id] = true
			if tc := diffTag(aTag, bTag); tc != nil {
				changes = append(changes, *tc)
			}
		} else {
			// new A Tag
			tc := wholeTagChange(aTag, ChangeTypeAdd)
			changes = append(changes, tc)
		}
	}

	for id, bTag := range bTags {
		if _, ok := existingMap[id]; ok {
			continue // already been diffed
		}

		tc := wholeTagChange(bTag, ChangeTypeRemove)
		changes = append(changes, tc)
	}

	return changes
}

func diffTag(a types.Tag, b types.Tag) *BaseChange {
	cc := &BaseChange{Id: a.Id, Name: a.Name, Type: ChangeTypeModify}

	fields := []fieldDiff{
		{"name", a.Name, b.Name},
		{"documentation", a.Documentation, b.Documentation},
		{"kind", a.Kind, b.Kind},
		{"value", a.Value, b.Value},
	}

	cc.Changes = diffFields(fields)

	if len(cc.Changes) == 0 {
		return nil
	}

	return cc
}

func wholeTagChange(c types.Tag, changeType ChangeType) BaseChange {
	cc := BaseChange{Id: c.Id, Name: c.Name, Type: changeType}

	fields := []fieldDiff{
		{"documentation", c.Documentation, ""},
		{"value", c.Value, ""},
		{"kind", c.Kind, ""},
	}

	cc.Changes = diffFieldsWithType(fields, changeType)

	return cc
}

package diff

import (
	"fmt"
	"mdj-diff/types"
)

func diffTags(a []types.Tag, b []types.Tag) []TagChange {
	var changes []TagChange

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

func diffTag(a types.Tag, b types.Tag) *TagChange {
	cc := &TagChange{Id: a.Id, Name: a.Name, Type: ChangeTypeModify}

	fields := []struct {
		name string
		a, b interface{}
	}{
		{"name", a.Name, b.Name},
		{"documentation", a.Documentation, b.Documentation},
		{"kind", a.Kind, b.Kind},
		{"value", a.Value, b.Value},
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

	if len(cc.Changes) == 0 {
		return nil
	}

	return cc
}

func wholeTagChange(c types.Tag, changeType ChangeType) TagChange {
	cc := TagChange{Id: c.Id, Name: c.Name, Type: changeType}

	// optional column fields
	if c.Documentation != "" {
		cc.Changes = append(cc.Changes, Change{Name: "documentation", Type: changeType, Value: c.Documentation})
	}

	if c.Value != "" {
		cc.Changes = append(cc.Changes, Change{Name: "value", Type: changeType, Value: c.Value})
	}

	if c.Kind != "" {
		cc.Changes = append(cc.Changes, Change{Name: "kind", Type: changeType, Value: c.Kind})
	}

	return cc
}

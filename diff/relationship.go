package diff

import (
	"mdj-diff/types"
)

func diffRelationships(a types.Entity, b types.Entity) []RelationshipChange {
	var changes []RelationshipChange

	// get column map
	aRels := a.GetRelationshipMap()
	bRels := b.GetRelationshipMap()

	existingMap := make(map[string]bool)

	for id, aRel := range aRels {
		if bRel, ok := bRels[id]; ok {
			existingMap[id] = true
			if r := diffRelationship(aRel, bRel); r != nil {
				changes = append(changes, *r)
			}
		} else {
			// new A Relationhsip
			r := wholeRelationshipChange(aRel, ChangeTypeAdd)
			changes = append(changes, r)
		}
	}

	for id, bRel := range bRels {
		if _, ok := existingMap[id]; ok {
			continue // already been diffed
		}

		r := wholeRelationshipChange(bRel, ChangeTypeRemove)
		changes = append(changes, r)
	}

	return changes
}

func diffRelationship(a types.Relationship, b types.Relationship) *RelationshipChange {
	r := &RelationshipChange{BaseChange: BaseChange{Name: a.End2.Reference.Ref, Type: ChangeTypeModify}}

	fields := []fieldDiff{
		{"documentation", a.Documentation, b.Documentation},
		{"end1.cardinality", a.End1.GetCardinality(), b.End1.GetCardinality()},
		{"end2.cardinality", a.End2.GetCardinality(), b.End2.GetCardinality()},
		{"end1.reference", a.End1.Reference.Ref, b.End1.Reference.Ref},
		{"end2.reference", a.End2.Reference.Ref, b.End2.Reference.Ref},
	}

	r.Changes = diffFields(fields)
	r.Tags = diffTags(a.Tags, b.Tags)

	if len(r.Changes)+len(r.Tags) == 0 {
		return nil
	}

	return r
}

func wholeRelationshipChange(c types.Relationship, changeType ChangeType) RelationshipChange {
	r := RelationshipChange{BaseChange: BaseChange{Name: c.End2.Reference.Ref, Type: changeType}}

	r.End1Cardinality = &Change{Name: "end1.cardinality", Type: changeType, Value: c.End1.GetCardinality()}
	r.End2Cardinality = &Change{Name: "end2.cardinality", Type: changeType, Value: c.End2.GetCardinality()}
	r.End1Reference = &Change{Name: "end1.reference", Type: changeType, Value: c.End1.Reference.Ref}
	r.End2Reference = &Change{Name: "end2.reference", Type: changeType, Value: c.End2.Reference.Ref}

	// optional relationship fields
	if c.Documentation != "" {
		r.Changes = append(r.Changes, Change{Name: "documentation", Type: changeType, Value: c.Documentation})
	}

	return r
}

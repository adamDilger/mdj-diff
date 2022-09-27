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
	r := &RelationshipChange{Name: a.End2.Reference.Ref, Type: ChangeTypeModify}

	change := false

	if a.End1.GetCardinality() != b.End1.GetCardinality() {
		change = true
		r.End1Cardinality = &Change{Name: "end1.cardinality", Type: ChangeTypeModify, Value: a.End1.GetCardinality(), Old: b.End1.GetCardinality()}
	}

	if a.End2.GetCardinality() != b.End2.GetCardinality() {
		change = true
		r.End2Cardinality = &Change{Name: "end2.cardinality", Type: ChangeTypeModify, Value: a.End2.GetCardinality(), Old: b.End2.GetCardinality()}
	}

	if a.End1.Reference.Ref != b.End1.Reference.Ref {
		change = true
		r.End1Reference = &Change{Name: "end1.reference", Type: ChangeTypeModify, Value: a.End1.Reference.Ref, Old: b.End1.Reference.Ref}
	}

	if a.End2.Reference.Ref != b.End2.Reference.Ref {
		change = true
		r.End2Reference = &Change{Name: "end2.reference", Type: ChangeTypeModify, Value: a.End2.Reference.Ref, Old: b.End2.Reference.Ref}
	}

	if a.GetDocumentation() != b.GetDocumentation() {
		change = true
		r.Changes = append(r.Changes, Change{Name: "documentation", Type: ChangeTypeModify, Value: a.GetDocumentation(), Old: b.GetDocumentation()})
	}

	r.Tags = diffTags(a.GetTags(), b.GetTags())

	// optional relationship fields
	if !change && len(r.Tags) == 0 {
		return nil
	}

	return r
}

func wholeRelationshipChange(c types.Relationship, changeType ChangeType) RelationshipChange {
	r := RelationshipChange{Name: c.End2.Reference.Ref, Type: changeType}

	r.End1Cardinality = &Change{Name: "end1.cardinality", Type: changeType, Value: c.End1.GetCardinality()}
	r.End2Cardinality = &Change{Name: "end2.cardinality", Type: changeType, Value: c.End2.GetCardinality()}
	r.End1Reference = &Change{Name: "end1.reference", Type: changeType, Value: c.End1.Reference.Ref}
	r.End2Reference = &Change{Name: "end2.reference", Type: changeType, Value: c.End2.Reference.Ref}

	// optional relationship fields
	if c.GetDocumentation() != "" {
		r.Changes = append(r.Changes, Change{Name: "documentation", Type: changeType, Value: c.GetDocumentation()})
	}

	return r
}

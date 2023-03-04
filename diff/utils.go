package diff

import "fmt"

type fieldDiff struct {
	name string
	a, b interface{}
}

func diffFields(fields []fieldDiff) []Change {
	return diffFieldsWithType(fields, ChangeTypeModify)
}

func diffFieldsWithType(fields []fieldDiff, changeType ChangeType) []Change {
	changes := []Change{}

	for _, f := range fields {
		if f.a != f.b {
			changes = append(changes, Change{
				Name:  f.name,
				Type:  changeType,
				Value: fmt.Sprintf("%v", f.a),
				Old:   fmt.Sprintf("%v", f.b),
			})
		}
	}

	return changes
}

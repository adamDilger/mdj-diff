package diff

import (
	"fmt"
	"mdj-diff/types"

	"github.com/jedib0t/go-pretty/v6/text"
)

type TableChange struct {
	Id            string
	Type          ChangeType
	Name          string
	Columns       []ColumnChange
	Relationships []RelationshipChange
	Changes       []Change
	Tags          []TagChange

	DataModel *types.Node
	Diagram   *types.Node
}

type ColumnChange struct {
	Id      string
	Name    string
	Type    ChangeType
	Changes []Change
	Tags    []TagChange
}

type RelationshipChange struct {
	Id              string
	Name            string
	Type            ChangeType
	End1Cardinality *Change
	End2Cardinality *Change
	End1Reference   *Change
	End2Reference   *Change
	Changes         []Change
	Tags            []TagChange
}

type TagChange struct {
	Id      string
	Name    string
	Type    ChangeType
	Changes []Change
}

type Change struct {
	Name  string
	Type  ChangeType
	Value interface{}
	Old   interface{}
}

type RelationshipEndChange struct {
}

type ChangeType string

const (
	ChangeTypeAdd    ChangeType = "add"
	ChangeTypeRemove ChangeType = "remove"
	ChangeTypeModify ChangeType = "modify"
)

var COLOR_ENABLED = true

/* need to sort out all these colour things */
var colorMap = map[ChangeType]text.Color{
	ChangeTypeRemove: text.FgRed,
	ChangeTypeAdd:    text.FgGreen,
	ChangeTypeModify: text.FgYellow,
}

var labelMap = map[ChangeType]string{
	ChangeTypeRemove: "Removed",
	ChangeTypeAdd:    "Added",
	ChangeTypeModify: "Modified",
}

var charMap = map[ChangeType]string{
	ChangeTypeRemove: "−",
	ChangeTypeAdd:    "＋",
	ChangeTypeModify: "~",
}

func (c ChangeType) Label() string { return labelMap[c] }
func (c ChangeType) Char() string  { return charMap[c] }

func (c ChangeType) SetColor() {
	if !COLOR_ENABLED {
		return
	}

	fmt.Printf("%s", colorMap[c].EscapeSeq())
}
func (c ChangeType) Sprintf(format string, args ...interface{}) string {
	if !COLOR_ENABLED {
		return fmt.Sprintf(format, args...)
	}

	return colorMap[c].Sprintf(format, args...)
}

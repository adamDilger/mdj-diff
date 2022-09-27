package printer

import (
	"encoding/json"
	"mdj-diff/diff"
	"os"
)

func PrintJson(changes []diff.TableChange) {
	e := json.NewEncoder(os.Stdout)
	e.SetIndent("", "  ")
	e.Encode(changes)
}

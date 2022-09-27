package printer

import (
	"fmt"
	"io"
	"mdj-diff/diff"
	"mdj-diff/types"
	"os"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

func PrintText(A, B *types.Project, changes []diff.TableChange) {
	print(&textPrinter{}, A, B, changes)
}

func PrintMarkdown(A, B *types.Project, changes []diff.TableChange) {
	diff.COLOR_ENABLED = false
	print(&markdownPrinter{}, A, B, changes)
}

func print(p printer, A, B *types.Project, changes []diff.TableChange) {
	p.RenderTitle(changes)

	for _, change := range changes {
		p.RenderEntityHeader(change)
		printTableAttributes(p, change)

		fmt.Println()

		if len(change.Columns) > 0 {
			p.RenderChangesetHeader("Columns", len(change.Columns))
		}

		for _, column := range change.Columns {
			printColumnChanges(p, column)
			fmt.Println()
		}

		if len(change.Relationships) > 0 {
			p.RenderChangesetHeader("Relationships", len(change.Relationships))
		}

		for _, rel := range change.Relationships {
			printRelationshipChanges(p, A, B, rel)
			fmt.Println()
		}

		if len(change.Tags) > 0 {
			p.RenderChangesetHeader("Tags", len(change.Tags))
		}

		for _, tag := range change.Tags {
			renderChanges(p, tag.Name, tag.Changes)
			fmt.Println()
		}
	}
}

func printTableAttributes(p printer, change diff.TableChange) {
	if len(change.Changes) == 0 {
		return
	}

	p.RenderChangesetHeader("Attributes", len(change.Changes))
	renderChanges(p, "", change.Changes)
}

func printColumnChanges(p printer, column diff.ColumnChange) {
	renderChanges(p, column.Name, column.Changes)

	if len(column.Tags) > 0 {
		fmt.Println()
		p.RenderChangesetHeader(fmt.Sprintf("[%s] Tags", column.Name), len(column.Tags))
	}

	for _, t := range column.Tags {
		renderChanges(p, t.Name, t.Changes)
	}
}

func printRelationshipChanges(p printer, A, B *types.Project, rel diff.RelationshipChange) {
	name := getNameForRef(A, B, rel.Name, rel.Type)

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetTitle(rel.Type.Sprintf(strings.ToUpper(name)))
	t.AppendHeader(table.Row{"", "Attribute", "Value", "Old"})
	t.SetStyle(table.StyleRounded)

	printRelationshipTable(t, A, B, rel.End1Reference)
	printRelationshipTable(t, A, B, rel.End2Reference)
	printRelationshipTable(t, A, B, rel.End1Cardinality)
	printRelationshipTable(t, A, B, rel.End2Cardinality)

	p.RenderTable(t)
}

func printRelationshipTable(t table.Writer, A, B *types.Project, c *diff.Change) {
	if c == nil {
		return
	}

	var value, old string
	if v, ok := c.Value.(string); ok {
		value = getNameForRef(A, B, v, c.Type)
	}
	if v, ok := c.Old.(string); ok {
		old = getNameForRef(A, B, v, c.Type)
	}

	p := c.Type.Sprintf
	row := table.Row{p(c.Type.Char()), p(c.Name), p("%v", value), p("%v", old)}
	t.AppendRow(row)
}

func renderChanges(printer printer, name string, changes []diff.Change) {
	renderChangesToWriter(os.Stdout, printer, name, changes)
}

func renderChangesToWriter(out io.Writer, printer printer, name string, changes []diff.Change) {
	if len(changes) == 0 {
		return
	}

	t := table.NewWriter()
	t.SetOutputMirror(out)
	if name != "" {
		t.SetTitle(name)
	}
	t.AppendHeader(table.Row{"", "Attribute", "Value", "Old"})
	t.SetStyle(table.StyleRounded)

	for _, c := range changes {
		p := c.Type.Sprintf
		var row table.Row

		switch c.Type {
		case diff.ChangeTypeRemove:
			row = table.Row{p(c.Type.Char()), p(c.Name), "", p("%v", c.Value)}
		case diff.ChangeTypeAdd:
			row = table.Row{p(c.Type.Char()), p(c.Name), p("%v", c.Value)}
		default:
			row = table.Row{p(c.Type.Char()), p(c.Name), p("%v", c.Value), text.FgRed.Sprintf("%v", c.Old)}
		}

		t.AppendRow(row)
	}

	printer.RenderTable(t)
}

func getNameForRef(A, B *types.Project, name string, t diff.ChangeType) string {
	var n types.Node

	switch t {
	case diff.ChangeTypeRemove:
		n = B.RefLookup[name]
	default:
		n = A.RefLookup[name]
	}

	if n == nil {
		return name
	}
	return n.GetName()
}

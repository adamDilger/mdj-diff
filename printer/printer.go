package printer

import (
	"fmt"
	"mdj-diff/diff"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

const WIDTH = 80

var EQUALS_LINE = strings.Repeat("=", WIDTH)

type printer interface {
	RenderTitle([]diff.TableChange)

	RenderChangesetHeader(name string, changeCount int)
	RenderEntityHeader(diff.TableChange)
	RenderTable(table.Writer)
}

// text

var _ printer = &textPrinter{}

type textPrinter struct{}

func (t *textPrinter) RenderTitle(changes []diff.TableChange) {
	fmt.Printf("%s\nMDJ Diff\n%s\n\n", EQUALS_LINE, EQUALS_LINE)
	fmt.Printf("Entities changed: %d\n\n", len(changes))
}

func (t *textPrinter) RenderEntityHeader(change diff.TableChange) {
	change.Type.SetColor()
	fmt.Println(EQUALS_LINE)
	fmt.Printf("%-*s%*s\n", WIDTH/2, change.Name, WIDTH/2, change.Type.Label())
	fmt.Println(EQUALS_LINE)
	fmt.Printf(text.EscapeReset)
}

func (t *textPrinter) RenderChangesetHeader(name string, changeCount int) {
	fmt.Printf("%s %d:\n", name, changeCount)
}
func (t *textPrinter) RenderTable(tw table.Writer) { tw.Render() }

// markdown

var _ printer = &markdownPrinter{}

type markdownPrinter struct{}

func (m *markdownPrinter) RenderEntityHeader(change diff.TableChange) {
	fmt.Printf("# ***%s*** (%s)\n", change.Name, change.Type.Label())
}

func (m *markdownPrinter) RenderTitle(changes []diff.TableChange) {
	fmt.Printf("# MDJ Diff\n\n")
	fmt.Printf("*Entities changed: %d*\n\n", len(changes))
}

func (m *markdownPrinter) RenderChangesetHeader(name string, changeCount int) {
	fmt.Printf("### %s (%d)\n", name, changeCount)
}

func (m *markdownPrinter) RenderTable(tw table.Writer) { tw.RenderMarkdown() }

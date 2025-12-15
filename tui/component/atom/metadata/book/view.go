package book

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/lrstanley/bubblezone"
)

func (m *M) View() string {
	var b strings.Builder

	bold := lipgloss.NewStyle().Bold(true)

	tab := table.New().Border(lipgloss.HiddenBorder()).Width(m.width)
	if m.hover {
		tab = tab.Border(lipgloss.NormalBorder())
	}
	tab = tab.Rows(
		[]string{bold.Render("Type"), R{}.BookType(m.book)},
		[]string{"", strings.Join(R{}.Genres(m.book), ", ")},
		[]string{bold.Render("Authors"), strings.Join(R{}.Authors(m.book), ", ")},
		[]string{bold.Render("Illustrators"), strings.Join(R{}.Illustrators(m.book), ", ")},
		[]string{bold.Render("Last Updated"), R{}.LastUpdated(m.book)},
	)

	b.WriteString(tab.String())
	return zone.Mark("TEST", b.String())
}

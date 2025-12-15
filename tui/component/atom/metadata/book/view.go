package book

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m *M) View() string {
	k := m.layout.Style().Bold(true)
	v := m.layout.WithMargin(m.layout.Margin + 1).Style().Foreground(lipgloss.Color("5"))
	return lipgloss.JoinVertical(
		lipgloss.Left,
		k.Foreground(lipgloss.Color("5")).Render(R{}.BookType(m.book)),
		k.Render("Genres"),
		v.Render(strings.Join(R{}.Genres(m.book), ", ")),
		k.Render("Authors"),
		v.Render(strings.Join(R{}.Authors(m.book), ", ")),
		k.Render("Illustrators"),
		v.Render(strings.Join(R{}.Illustrators(m.book), ", ")),
		k.Render("Last Updated"),
		v.Render(R{}.LastUpdated(m.book)),
	)
}

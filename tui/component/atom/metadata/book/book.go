package book

import (
	"strings"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/minkezhang/truffle-api/db/atom/metadata/book"

	model_ui "github.com/minkezhang/truffle/tui/component/util/model"
)

type O struct {
	model_ui.O

	Book *book.M
}

func Make(o O) M {
	return M{
		Base: model_ui.Make(o.O),
		book: o.Book,
	}
}

type M struct {
	model_ui.Base

	book *book.M
}

func (m M) Init() tea.Cmd                           { return nil }
func (m M) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return m, nil }

func (m M) View() string {
	k := m.Column().Style().Bold(true)
	v := m.Column().WithMargin(m.Column().Margin + 1).Style().Foreground(lipgloss.Color("5"))
	return m.RenderOrDie(
		lipgloss.JoinVertical(
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
		),
	)
}

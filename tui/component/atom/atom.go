package atom

import (
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/minkezhang/truffle-api/db/atom"
	"github.com/minkezhang/truffle-api/db/atom/metadata/book"

	image_ui "github.com/minkezhang/truffle/tui/component/image"
	book_ui "github.com/minkezhang/truffle/tui/component/metadata/book"
)

type O struct {
	Atom *atom.A
}

type M struct {
	atom     *atom.A
	metadata *book_ui.M
	image    *image_ui.M
}

func New(o O) *M {
	return &M{
		atom: o.Atom,
		metadata: book_ui.New(book_ui.O{
			Book: o.Atom.Metadata().(*book.M),
		}),
		image: image_ui.New(image_ui.O{
			Width: 30,
			URL:   o.Atom.PreviewURL(),
		}),
	}
}

func (m *M) Init() tea.Cmd {
	return tea.Batch(
		m.image.Init(),
		m.metadata.Init(),
	)
}

func (m *M) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	for _, n := range []tea.Model{
		m.image,
		m.metadata,
	} {
		_, c := n.Update(msg)
		cmds = append(cmds, c)
	}

	return m, tea.Batch(cmds...)
}

func (m *M) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			m.image.View(),
			m.metadata.View(),
		),
		lipgloss.NewStyle().Width(80).Render(m.atom.Synopsis()),
	)
}

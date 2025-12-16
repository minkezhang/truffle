package atom

import (
	"fmt"
	"math"
	"strings"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/minkezhang/truffle-api/db/atom"
	"github.com/minkezhang/truffle-api/db/atom/metadata/book"
	"github.com/minkezhang/truffle/tui/util/grid"

	epb "github.com/minkezhang/truffle-api/proto/go/enums"
	image_ui "github.com/minkezhang/truffle/tui/component/atom/image"
	book_ui "github.com/minkezhang/truffle/tui/component/atom/metadata/book"
	score_ui "github.com/minkezhang/truffle/tui/component/atom/score"
	titles_ui "github.com/minkezhang/truffle/tui/component/atom/titles"
)

type O struct {
	Layout         grid.L // Total width of the element
	Atom           *atom.A
	CacheDirectory string
}

type M struct {
	layout grid.L
	grid   grid.G

	atom     *atom.A
	metadata tea.Model
	image    tea.Model
	score    tea.Model
	titles   tea.Model
}

func New(o O) *M {
	l := o.Layout
	g := l.Grid(3)

	return &M{
		layout: l,
		grid:   g,
		atom:   o.Atom,
		metadata: book_ui.New(book_ui.O{
			Book:   o.Atom.Metadata().(*book.M),
			Layout: g.Column(2, 1, 0),
		}),
		image: image_ui.New(image_ui.O{
			CacheDirectory: o.CacheDirectory,
			Layout:         g.Column(1, 0, 0),
			Height:         height(g.Column(1, 0, 0).Content),
			URL:            o.Atom.PreviewURL(),
		}),
		score: score_ui.New(score_ui.O{
			Score: o.Atom.Score(),
		}),
		titles: titles_ui.New(titles_ui.O{
			Titles: o.Atom.Titles(),
		}),
	}
}

func (m *M) Init() tea.Cmd {
	return tea.Batch(
		m.image.Init(),
		m.metadata.Init(),
		m.score.Init(),
		m.titles.Init(),
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

func height(width int) int {
	return int( // A4 ratio; in pixels
		math.Trunc(float64(width) * 1.414),
	)
}

func (m *M) View() string {
	var id string
	switch t := m.atom.APIType(); t {
	case epb.API_API_VIRTUAL:
		id = "truffle"
	default:
		id = fmt.Sprintf(
			"%s/%s",
			strings.ToLower(strings.ReplaceAll(m.atom.APIType().String(), "API_", "")),
			m.atom.APIID(),
		)
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		m.grid.Column(3, 0, 0).Style().MarginBottom(1).Render(m.titles.View()),
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			m.image.View(),
			lipgloss.JoinVertical(
				lipgloss.Left,
				m.grid.Column(2, 1, 0).Style().Bold(true).Foreground(lipgloss.Color("5")).Render(id),
				m.metadata.View(),
				m.grid.Column(2, 1, 0).Style().MarginTop(1).Render(m.score.View()),
			),
		),
		m.grid.Column(3, 0, 0).Style().Render(m.atom.Synopsis()),
	)
}

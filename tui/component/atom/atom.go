package atom

import (
	"fmt"
	"math"
	"strings"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/minkezhang/truffle-api/db/atom"
	"github.com/minkezhang/truffle-api/db/atom/metadata/book"

	image_ui "github.com/minkezhang/truffle/tui/component/atom/image"
	book_ui "github.com/minkezhang/truffle/tui/component/atom/metadata/book"
	score_ui "github.com/minkezhang/truffle/tui/component/atom/score"
	titles_ui "github.com/minkezhang/truffle/tui/component/atom/titles"
)

type O struct {
	Width          int // Total width of the element
	Atom           *atom.A
	CacheDirectory string
}

type M struct {
	width    int
	atom     *atom.A
	metadata tea.Model
	image    tea.Model
	score    tea.Model
	titles   tea.Model
}

func New(o O) *M {
	return &M{
		width: o.Width,
		atom:  o.Atom,
		metadata: book_ui.New(book_ui.O{
			Book:  o.Atom.Metadata().(*book.M),
			Width: column(o.Width) * 2,
		}),
		image: image_ui.New(image_ui.O{
			CacheDirectory: o.CacheDirectory,
			Width:          column(o.Width),
			Height:         height(column(o.Width)),
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

func column(width int) int { return width / 3 }

func height(width int) int {
	return int( // A4 ratio; in pixels
		math.Trunc(float64(width) * 1.414),
	)
}

func (m *M) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			m.image.View(),
			lipgloss.JoinVertical(
				lipgloss.Left,
				lipgloss.NewStyle().Width(m.width).Bold(true).Foreground(lipgloss.Color("5")).Render(
					fmt.Sprintf(
						"%s/%s",
						strings.ToLower(strings.ReplaceAll(m.atom.APIType().String(), "API_", "")),
						m.atom.APIID(),
					),
				),
				m.metadata.View(),
			),
		),
		lipgloss.NewStyle().Width(m.width).Render(m.titles.View()),
		lipgloss.NewStyle().Width(m.width).Margin(1, 0, 1, 0).Render(m.score.View()),
		lipgloss.NewStyle().Width(m.width).Render(m.atom.Synopsis()),
	)
}

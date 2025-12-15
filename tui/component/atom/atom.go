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

var (
	width  int = 30                                      // In pixels
	height int = int(math.Trunc(float64(width) * 1.414)) // A4 ratio; in pixels
)

type O struct {
	Atom           *atom.A
	CacheDirectory string
}

type M struct {
	atom     *atom.A
	metadata *book_ui.M
	image    *image_ui.M
	score    *score_ui.M
	titles   *titles_ui.M
}

func New(o O) *M {
	return &M{
		atom: o.Atom,
		metadata: book_ui.New(book_ui.O{
			Book:  o.Atom.Metadata().(*book.M),
			Width: width * 2,
		}),
		image: image_ui.New(image_ui.O{
			CacheDirectory: o.CacheDirectory,
			Width:          width,
			Height:         height,
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

func (m *M) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			lipgloss.JoinVertical(
				lipgloss.Left,
				lipgloss.Place(
					width,
					// Two vertical pixels per character may leave
					// a pixel unaccounted for.
					height/2+1,
					lipgloss.Top,
					lipgloss.Center,
					m.image.View(),
				),
			),
			lipgloss.JoinVertical(
				lipgloss.Left,
				m.metadata.View(),
				fmt.Sprintf(
					"%s:%s",
					strings.ReplaceAll(m.atom.APIType().String(), "API_", ""),
					m.atom.APIID(),
				),
				m.score.View(),
			),
		),
		lipgloss.NewStyle().Width(width*3).Render(m.titles.View()),
		lipgloss.NewStyle().Width(width*3).Render(m.atom.Synopsis()),
	)
}

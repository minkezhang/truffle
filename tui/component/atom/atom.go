package atom

import (
	"fmt"
	"reflect"
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
	model_ui "github.com/minkezhang/truffle/tui/component/util/model"
)

type O struct {
	model_ui.O
	Atom           *atom.A
	CacheDirectory string
}

type M struct {
	*model_ui.Base
	grid grid.G

	atom     *atom.A
	metadata tea.Model
	image    tea.Model
	score    tea.Model
	titles   tea.Model
}

func New(o O) *M {
	c := o.Column
	g := c.Grid(3)

	return &M{
		Base: model_ui.New(o.O),
		grid: g,
		atom: o.Atom,
		metadata: book_ui.New(book_ui.O{
			O: model_ui.O{
				Column: g.Column(2).WithMargin(1),
			},
			Book: o.Atom.Metadata().(*book.M),
		}),
		image: image_ui.New(image_ui.O{
			O: model_ui.O{
				Column: g.Column(1),
			},
			CacheDirectory: o.CacheDirectory,
			URL:            o.Atom.PreviewURL(),
		}),
		score: score_ui.New(score_ui.O{
			O: model_ui.O{
				Column: g.Column(2),
			},
			Score: o.Atom.Score(),
		}),
		titles: titles_ui.New(titles_ui.O{
			O: model_ui.O{
				Column: g.Column(3),
			},
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

type updateAtomMsg struct {
	model_ui.BaseMsg
	payload *atom.A
}

func UpdateAtomAsync(m tea.Model, v *atom.A) tea.Cmd {
	if m, ok := m.(*M); ok {
		return func() tea.Msg {
			return updateAtomMsg{
				BaseMsg: model_ui.BaseMsg{
					ID: m.ID(),
				},
				payload: v,
			}
		}
	}
	return model_ui.ErrorCmd(fmt.Errorf("incorrect model type: %v", reflect.TypeOf(m)))
}

func (m *M) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case updateAtomMsg:
		if m.ID() == msg.ID {
			m.atom = msg.payload
			cmds = append(cmds, tea.Batch(
				image_ui.UpdateImageAsync(m.image, m.atom.PreviewURL()),
				book_ui.UpdateBookAsync(m.metadata, m.atom.Metadata().(*book.M)),
				score_ui.UpdateScoreAsync(m.score, m.atom.Score()),
				titles_ui.UpdateTitlesAsync(m.titles, m.atom.Titles()),
			))
		}
	}

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

	return m.RenderOrDie(
		lipgloss.JoinVertical(
			lipgloss.Left,
			m.grid.Column(3).Style().MarginBottom(1).Render(m.titles.View()),
			lipgloss.JoinHorizontal(
				lipgloss.Top,
				m.image.View(),
				lipgloss.JoinVertical(
					lipgloss.Left,
					m.grid.Column(2).WithMargin(1).Style().Bold(true).Foreground(lipgloss.Color("5")).Render(id),
					m.metadata.View(),
					m.grid.Column(2).WithMargin(1).Style().MarginTop(1).Render(m.score.View()),
				),
			),
			m.grid.Column(3).Style().Render(m.atom.Synopsis()),
		),
	)
}

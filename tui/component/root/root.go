package root

import (
	"context"

	"github.com/charmbracelet/bubbletea"
	"github.com/lrstanley/bubblezone"
	"github.com/minkezhang/truffle-api/client/mal"
	"github.com/minkezhang/truffle-api/client/query"

	epb "github.com/minkezhang/truffle-api/proto/go/enums"
	atom_ui "github.com/minkezhang/truffle/tui/component/atom"
)

type O struct {
	CacheDirectory string
}

type M struct {
	directory string
	atom      tea.Model
}

func New(o O) *M {
	c := mal.New(mal.O{
		ClientID:         "6114d00ca681b7701d1e15fe11a4987e",
		PopularityCutoff: 10000,
		MaxResults:       2,
		NSFW:             true,
	})
	a, _ := c.Get(context.Background(), query.G{
		AtomType: epb.Type_TYPE_BOOK,
		ID:       "107562",
	})

	return &M{
		directory: o.CacheDirectory,
		atom: atom_ui.New(atom_ui.O{
			Width:          90,
			CacheDirectory: o.CacheDirectory,
			Atom:           a,
		}),
	}
}

func (m *M) Init() tea.Cmd { return m.atom.Init() }

func (m *M) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			fallthrough
		case tea.KeyCtrlD:
			return m, tea.Quit
		case tea.KeyCtrlZ:
			return m, tea.Suspend
		}
	case tea.WindowSizeMsg: // Clear buffer
		return m, tea.ClearScreen
	}
	_, c := m.atom.Update(msg)
	return m, c
}

func (m *M) View() string {
	return zone.Scan(m.atom.View())
}

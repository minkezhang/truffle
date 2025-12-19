package node

import (
	"fmt"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/minkezhang/truffle-api/db/atom"
	"github.com/minkezhang/truffle-api/db/node"
	"github.com/minkezhang/truffle/tui/util/grid"

	atom_ui "github.com/minkezhang/truffle/tui/component/atom"
	source_list_ui "github.com/minkezhang/truffle/tui/component/node/source_list"
	model_ui "github.com/minkezhang/truffle/tui/component/util/model"
)

type O struct {
	model_ui.O

	CacheDirectory string
	Node           *node.N
}

type M struct {
	*model_ui.Base
	grid grid.G

	directory string
	node      *node.N

	sources tea.Model
	atoms   map[source_list_ui.V]tea.Model
	focus   source_list_ui.V
}

func Make(o O) M {
	g := o.Column.Grid(4)
	m := M{
		Base:      model_ui.New(o.O),
		grid:      g,
		directory: o.CacheDirectory,
		node:      o.Node,
		atoms:     map[source_list_ui.V]tea.Model{},
	}
	vs := []source_list_ui.V{}
	for _, a := range append([]*atom.A{
		o.Node.Virtual(),
	}, o.Node.Atoms()...) {
		v := source_list_ui.V{
			API: a.APIType(),
			ID:  a.APIID(),
		}
		m.atoms[v] = atom_ui.Make(atom_ui.O{
			O: model_ui.O{
				Column: g.Column(3),
			},
			CacheDirectory: o.CacheDirectory,
			Atom:           a,
		})
		vs = append(vs, v)
	}

	if len(vs) > 0 {
		m.focus = vs[0]
	}

	m.sources = source_list_ui.New(source_list_ui.O{
		O: model_ui.O{
			Column: m.grid.Column(3),
		},
		Values: vs,
	})

	return m
}

func (m M) Init() tea.Cmd {
	cmds := []tea.Cmd{}

	for _, a := range m.atoms {
		cmds = append(cmds, a.Init())
	}
	cmds = append(cmds, m.sources.Init())

	return tea.Batch(cmds...)
}

func (m M) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case source_list_ui.SelectMsg:
		m.focus = msg.Focus
	}

	for _, a := range m.atoms {
		_, c := a.Update(msg)
		cmds = append(cmds, c)
	}
	_, c := m.sources.Update(msg)
	cmds = append(cmds, c)

	return m, tea.Batch(cmds...)
}

func (m M) View() string {
	var notes = lipgloss.NewStyle().Foreground(lipgloss.Color("5")).Render(m.node.Notes())
	if m.node.Notes() == "" {
		notes = lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Render("N/A")
	}

	atom := ""
	if a, ok := m.atoms[m.focus]; ok {
		atom = a.View()
	}

	return m.RenderOrDie(
		lipgloss.JoinVertical(
			lipgloss.Right,
			m.sources.View(),

			lipgloss.JoinHorizontal(
				lipgloss.Top,
				lipgloss.JoinVertical(
					lipgloss.Left,
					m.grid.Column(1).WithMargin(1).Style().Bold(true).Render("Queued"),
					m.grid.Column(1).WithMargin(1).WithPadding(1).Style().Foreground(lipgloss.Color("5")).Render(fmt.Sprintf("%v", m.node.IsQueued())),
					m.grid.Column(1).WithMargin(1).Style().Bold(true).Render("Notes"),
					m.grid.Column(1).WithMargin(1).Style().Render(notes),
				),
				atom,
			),
		),
	)
}

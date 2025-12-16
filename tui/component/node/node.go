package node

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/minkezhang/truffle-api/db/atom"
	"github.com/minkezhang/truffle-api/db/node"
	"github.com/minkezhang/truffle/tui/util/grid"

	epb "github.com/minkezhang/truffle-api/proto/go/enums"
	atom_ui "github.com/minkezhang/truffle/tui/component/atom"
	tab_ui "github.com/minkezhang/truffle/tui/component/node/tab"
)

type O struct {
	Layout         grid.L
	CacheDirectory string
	Node           *node.N
}

type M struct {
	layout    grid.L
	grid      grid.G
	directory string

	node *node.N

	atoms  []tea.Model
	aIndex int

	tabs []*tab_ui.M
}

func New(o O) *M {
	g := o.Layout.Grid(4)
	m := &M{
		layout:    o.Layout,
		grid:      g,
		directory: o.CacheDirectory,
		node:      o.Node,
	}
	for _, a := range append([]*atom.A{
		o.Node.Virtual(),
	}, o.Node.Atoms()...) {
		m.atoms = append(m.atoms, atom_ui.New(atom_ui.O{
			Layout:         g.Column(3, 0, 0),
			CacheDirectory: o.CacheDirectory,
			Atom:           a,
		}))

		var v string
		switch t := a.APIType(); t {
		case epb.API_API_VIRTUAL:
			v = "TRUFFLE"
		default:
			v = fmt.Sprintf(
				"%s/%s",
				strings.ReplaceAll(a.APIType().String(), "API_", ""),
				a.APIID(),
			)
		}

		m.tabs = append(m.tabs, tab_ui.New(tab_ui.O{
			Value: v,
			State: tab_ui.StateNone,
		}))
	}
	return m
}

func (m *M) Init() tea.Cmd {
	cmds := []tea.Cmd{}
	for _, a := range m.atoms {
		cmds = append(cmds, a.Init())
	}
	for _, t := range m.tabs {
		cmds = append(cmds, t.Init())
	}
	cmds = append(cmds, func() tea.Msg {
		return tab_ui.FocusMsg{
			Key:    m.tabs[m.aIndex].Key(),
			Target: tab_ui.StateActive,
		}
	})
	return tea.Batch(cmds...)
}

func (m *M) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	for _, n := range append([]tea.Model{}, m.atoms...) {
		_, c := n.Update(msg)
		cmds = append(cmds, c)
	}
	for _, t := range append([]*tab_ui.M{}, m.tabs...) {
		_, c := t.Update(msg)
		cmds = append(cmds, c)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyLeft:
			src := m.aIndex
			m.aIndex = (m.aIndex - 1) % len(m.atoms)
			if m.aIndex < 0 {
				m.aIndex = len(m.atoms) - 1
			}
			dst := m.aIndex

			cmds = append(cmds, tea.Sequence(
				func() tea.Msg {
					return tab_ui.FocusMsg{
						Key:    m.tabs[src].Key(),
						Target: tab_ui.StateNone,
					}
				},
				func() tea.Msg {
					return tab_ui.FocusMsg{
						Key:    m.tabs[dst].Key(),
						Target: tab_ui.StateActive,
					}
				},
			))
		case tea.KeyRight:
			src := m.aIndex
			m.aIndex = (m.aIndex + 1) % len(m.atoms)
			dst := m.aIndex

			cmds = append(cmds, tea.Sequence(
				func() tea.Msg {
					return tab_ui.FocusMsg{
						Key:    m.tabs[src].Key(),
						Target: tab_ui.StateNone,
					}
				},
				func() tea.Msg {
					return tab_ui.FocusMsg{
						Key:    m.tabs[dst].Key(),
						Target: tab_ui.StateActive,
					}
				},
			))
		}
	}

	return m, tea.Batch(cmds...)
}

func (m *M) View() string {
	var tabs []string
	for _, t := range m.tabs {
		tabs = append(tabs, t.View())
	}

	var notes = lipgloss.NewStyle().Render(m.node.Notes())
	if notes == "" {
		notes = lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Render("N/A")
	}

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		lipgloss.JoinVertical(
			lipgloss.Left,
			lipgloss.JoinHorizontal(
				lipgloss.Top,
				tabs...,
			),
			m.atoms[m.aIndex].View(),
		),
		lipgloss.JoinVertical(
			lipgloss.Left,
			m.grid.Column(1, 1, 0).Style().Bold(true).Render("Queued"),
			m.grid.Column(1, 1, 1).Style().Foreground(lipgloss.Color("5")).Render(fmt.Sprintf("%v", m.node.IsQueued())),
			m.grid.Column(1, 1, 0).Style().MarginTop(1).MarginBottom(1).Bold(true).Render("Notes"),
			m.grid.Column(1, 1, 0).Style().Render(notes),
		),
	)
}

package directory

import (
	"github.com/charmbracelet/bubbletea"
)

type Node interface {
	tea.Model

	ID() string
	Parent() string

	Focus(i int) tea.Cmd
	Blur() tea.Cmd
	CurrentIndex() int
}

type TabMsg struct {
	ID           string
	CurrentIndex int // Current index
}

type RegisterNodeMsg struct {
	Node Node
}

type D struct {
	nodes    map[string]Node
	children map[string][]string
	parent   map[string]string

	current_node string
}

func New() *D {
	return &D{
		nodes:    map[string]Node{},
		children: map[string][]string{},
		parent:   map[string]string{},
	}
}

func (d *D) Init() tea.Cmd { return nil }
func (d *D) View() string  { return "" }

func (d *D) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case RegisterNodeMsg:
		d.nodes[msg.Node.ID()] = msg.Node
		d.children[msg.Node.Parent()] = append(d.children[msg.Node.Parent()], msg.Node.ID())
		d.parent[msg.Node.Parent()] = msg.Node.ID()
	case TabMsg:
		cmds = append(
			cmds,
			d.nodes[d.current_node].Blur(),
			d.nodes[msg.ID].Focus(d.nodes[msg.ID].CurrentIndex()+1),
		)
		d.current_node = msg.ID
	}
	return d, tea.Batch(cmds...)
}

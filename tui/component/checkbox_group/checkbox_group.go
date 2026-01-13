package checkbox_group

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbletea"
	"github.com/minkezhang/truffle/tui/component/checkbox"
	"github.com/minkezhang/truffle/tui/component/directory/base"
	"github.com/minkezhang/truffle/tui/component/focusable"
)

type Node struct {
	*focusable.Node

	is_radio bool
	values   []*checkbox.Node
}

type Value struct {
	Label      string
	Value      string
	IsSelected bool
}

type O struct {
	Prefix   string
	ParentID string
	IsRadio  bool
	Values   []Value
}

func New(o O) *Node {
	n := &Node{
		Node:     focusable.New(o.Prefix, o.ParentID, 0),
		is_radio: o.IsRadio,
		values:   []*checkbox.Node{},
	}
	for _, v := range o.Values {
		b := checkbox.New(checkbox.O{
			Prefix:     fmt.Sprintf("%s-%s", o.Prefix, v.Value),
			ParentID:   n.ID(),
			Label:      v.Label,
			Value:      v.Value,
			IsSelected: v.IsSelected,
			IsRadio:    o.IsRadio,
		})
		n.values = append(n.values, b)
	}
	return n
}

func (n *Node) Init() tea.Cmd {
	var cmds []tea.Cmd
	for _, b := range n.values {
		cmds = append(cmds, b.Init())
	}
	return tea.Sequence(
		tea.Sequence(
			cmds..., // preserve tab order
		),
		func() tea.Msg {
			return base.RegisterMessage{
				Node: n,
			}
		},
	)
}

func (n *Node) Selected() []string {
	var vs []string
	for _, b := range n.values {
		if (*checkbox.Node)(b).IsSelected() {
			vs = append(vs, (*checkbox.Node)(b).Value())
		}
	}
	return vs
}

func (n *Node) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case checkbox.SelectCheckboxInput:
		// Radio buttons are exclusively selected.
		if msg.ParentID == n.ID() && n.is_radio && msg.IsSelected {
			for _, b := range n.values {
				if b.ID() != msg.ID {
					cmds = append(cmds, func() tea.Msg {
						return checkbox.SelectCheckboxInput{
							ID:         b.ID(),
							ParentID:   n.ID(),
							IsSelected: false,
						}
					})
				}
			}
		}
	}
	for _, b := range n.values {
		_, c := b.Update(msg)
		cmds = append(cmds, c)
	}
	return n, tea.Batch(cmds...)
}

func (n *Node) View() string {
	var parts []string
	for _, b := range n.values {
		parts = append(parts, b.View())
	}
	return strings.Join(parts, " ") // TODO(minkezhang): Render within bounds.
}

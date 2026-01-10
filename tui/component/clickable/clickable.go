// Package clickable implements an action node which tracks mouse clicks and
// ensures a signal is only sent if both the mouse down and mouse up events are
// within the same bounding box.
//
// Example:
//
//	type N struct {  // tea.Cmd
//	  cl *clickable.Node
//	}
//
//	func (n *N) Init() tea.Cmd { return tea.Batch(n.cl.Init(), ...) }
//
//	func (n *N) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
//	  switch msg := msg.(type) {
//	  case clickable.Click:
//	    if msg.ID == n.cl.ID() { ... }
//	  }
//	}
//
//	func (n *N) View() string { return zone.Mark(n.cl.ID(), ...) }
package clickable

import (
	"github.com/charmbracelet/bubbletea"
	"github.com/lrstanley/bubblezone"
	"github.com/minkezhang/truffle/tui/component/focusable/directory/types"
	"github.com/minkezhang/truffle/tui/component/focusable"
)

type Click struct {
	ID string
}

type Node struct {
	*focusable.Node

	is_down bool
}

func New(parent_id string) *Node {
	return &Node{
		Node: focusable.New("clickable", parent_id, 0),
	}
}

func (n *Node) Init() tea.Cmd {
	return func() tea.Msg {
		return types.RegisterNodeMessage{
			Node: n,
		}
	}
}

func (n *Node) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.MouseMsg:
		if msg.Button == tea.MouseButtonLeft {
			if msg.Action == tea.MouseActionPress {
				if zone.Get(n.ID()).InBounds(msg) {
					n.is_down = true
				}
			} else if msg.Action == tea.MouseActionRelease {
				if zone.Get(n.ID()).InBounds(msg) {
					if n.is_down {
						cmds = append(cmds, func() tea.Msg {
							return Click{ID: n.ID()}
						})
					}
				}
				n.is_down = false
			}
		}
	}
	return n, tea.Batch(cmds...)
}

func (n *Node) View() string { return "" }

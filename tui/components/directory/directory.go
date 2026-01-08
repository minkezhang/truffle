package directory

import (
	"fmt"
	"slices"
	"strings"

	"github.com/charmbracelet/bubbletea"
	"github.com/minkezhang/truffle/tui/components/directory/types"
)

type D struct {
	nodes    map[string]types.Node
	children map[string][]string
	parent   map[string]string

	current_node_id string
	dirty           bool
	order_cache     []string
}

func New() *D {
	return &D{
		nodes:    map[string]types.Node{},
		children: map[string][]string{},
		parent:   map[string]string{},
		dirty:    true,
	}
}

// order returns the pre-order traversal of the directory.
func (d *D) order() []string {
	if !d.dirty {
		return d.order_cache
	}
	var f func(n string) []string
	f = func(n string) []string {
		var result = []string{}
		if n != "" {
			result = append(result, n)
		}
		for _, c := range d.children[n] {
			result = append(result, f(c)...)
		}
		return result
	}
	d.order_cache = f("")
	return d.order_cache
}

func (d *D) Init() tea.Cmd { return nil }

func (d *D) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch t := msg.Type; t {
		case tea.KeyShiftTab:
			dst := d.nodes[d.current_node_id].FocusIndex() - 1
			cmds = append(cmds, d.nodes[d.current_node_id].OnFocus(dst))
		case tea.KeyTab:
			dst := d.nodes[d.current_node_id].FocusIndex() + 1
			cmds = append(cmds, d.nodes[d.current_node_id].OnFocus(dst))
		}
	case types.RegisterNodeMessage:
		d.nodes[msg.Node.ID()] = msg.Node
		d.children[msg.Node.ParentID()] = append(d.children[msg.Node.ParentID()], msg.Node.ID())
		d.parent[msg.Node.ID()] = msg.Node.ParentID()
		d.dirty = true
	case types.FocusMessage:
		if n, ok := d.nodes[d.current_node_id]; ok {
			cmds = append(cmds, n.OnBlur())
		}
		if n, ok := d.nodes[msg.ID]; ok {
			cmds = append(cmds, n.OnFocus(msg.Index))
		} else {
			// TODO(minkezhang): Raise error.
		}
		d.current_node_id = msg.ID
	case types.EOFMessage:
		target_index := slices.IndexFunc(d.order(), func(v string) bool { return v == d.current_node_id })
		if msg.IsHead {
			target_index = target_index - 1
			if target_index < 0 {
				target_index = len(d.order()) - 1
			}
		} else {
			target_index = target_index + 1
			if target_index >= len(d.order()) {
				target_index = 0
			}
		}
		target_id := d.order()[target_index]
		focus_index := 0
		if msg.IsHead {
			focus_index = d.nodes[target_id].NElements() - 1
		} else {
			focus_index = 0
		}
		cmds = append(
			cmds,
			d.nodes[d.current_node_id].OnBlur(),
			d.nodes[target_id].OnFocus(focus_index),
		)
		d.current_node_id = target_id
	}
	return d, tea.Batch(cmds...)
}

func (d *D) View() string {
	var tree func(indent int, prefix string, nodes []string) []string
	// From github.com/campoy/tools/tree.
	tree = func(indent int, prefix string, nodes []string) []string {
		result := []string{}
		for i, n := range nodes {
			directory := "│  "
			file := "├─ "
			if i == len(nodes) - 1 {
				directory = "   "
				file = "└─ "
			}
			if n == "" {
				directory = " "
				result = append(result, "(root)")
			} else {
				result = append(
					result,
					fmt.Sprintf(
						"%v%v%v",
						prefix,
						file,
						n,
					),
				)
			}
			result = append(
				result,
				tree(indent+1, prefix + directory, d.children[n])...,
			)
		}
		return result
	}
	return strings.Join(tree(0, "", []string{""}), "\n")
}

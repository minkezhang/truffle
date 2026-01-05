package directory

import (
	"slices"

	"github.com/charmbracelet/bubbletea"
)

type FocusState int

const (
	FocusStateNone FocusState = iota
	FocusStateActive
)

type FocusHandler interface {
	OnFocus(i int) tea.Cmd
	OnBlur() tea.Cmd
	NElements() int
}

type Node interface {
	tea.Model

	ID() string
	Parent() string

	FocusHandler
	CurrentIndex() int
	FocusState() FocusState
}

type RegisterNode struct {
	Node Node
}

type EOF struct {
	ID     string
	IsHead bool
}

type Focus struct {
	ID    string
	Index int
}

type D struct {
	nodes    map[string]Node
	children map[string][]string
	parent   map[string]string

	current_node string
	dirty       bool
	order_cache []string
}

func New() *D {
	return &D{
		nodes:    map[string]Node{},
		children: map[string][]string{},
		parent:   map[string]string{},
		dirty:    true,
	}
}

// order returns the in-order traversal of the directory.
func (d *D) order() []string {
	if !d.dirty {
		return d.order_cache
	}
	var f func(n string) []string
	f = func(n string) []string {
		var result = []string{n}
		for _, c := range d.children[n] {
			result = append(result, f(c)...)
		}
		return result
	}
	d.order_cache = f("")
	return d.order_cache
}

func (d *D) Init() tea.Cmd { return nil }
func (d *D) View() string  { return "" }

func (d *D) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch t := msg.Type; t {
		case tea.KeyShiftTab:
			dst := d.nodes[d.current_node].CurrentIndex() - 1
			cmds = append(cmds, d.nodes[d.current_node].OnFocus(dst))
		case tea.KeyTab:
			dst := d.nodes[d.current_node].CurrentIndex() + 1
			cmds = append(cmds, d.nodes[d.current_node].OnFocus(dst))
		}
	case RegisterNode:
		d.nodes[msg.Node.ID()] = msg.Node
		d.children[msg.Node.Parent()] = append(d.children[msg.Node.Parent()], msg.Node.ID())
		d.parent[msg.Node.ID()] = msg.Node.Parent()
		d.dirty = true
	case Focus:
		cmds = append(
			cmds,
			d.nodes[d.current_node].OnBlur(),
			d.nodes[msg.ID].OnFocus(msg.Index),
		)
		d.current_node = msg.ID
	case EOF:
		target_index := slices.IndexFunc(d.order(), func(v string) bool { return v == d.current_node })
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
		target := d.order()[target_index]
		focus_index := 0
		if msg.IsHead {
			focus_index = d.nodes[target].NElements() - 1
		} else {
			focus_index = 0
		}
		cmds = append(
			cmds,
			d.nodes[d.current_node].OnBlur(),
			d.nodes[target].OnFocus(focus_index),
		)
		d.current_node = target
	}
	return d, tea.Batch(cmds...)
}

package focusable

import (
	"slices"

	"github.com/charmbracelet/bubbletea"
	"github.com/minkezhang/truffle/tui/component/directory/base"
	"github.com/minkezhang/truffle/tui/component/directory/focusable/types"
)

type Node interface {
	base.Identifiable

	// OnFocus may be called in the Update method by the directory. This
	// function directly instructs the node to activate any sub-elements for
	// input, e.g. focusing a specific input field.
	OnFocus(i int) tea.Cmd

	// OnBlur may be called in the Update method by the directory. This
	// function directly instructs the node to deactivate all sub-elements.
	OnBlur() tea.Cmd

	// NElements returns the number of interactable (and therefore
	// tab-focusable) elements in this node. NElements does not take into
	// account the number of [Node] children.
	NElements() int

	FocusIndex() int
	FocusState() types.FocusState
}

type D struct {
	directory *base.D

	current_node_id string
	dirty           bool
	order_cache     []string
}

func New() *D {
	return &D{
		directory: base.New(),
		dirty:     true,
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
		for _, c := range d.directory.Children[n] {
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
	var c tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch t := msg.Type; t {
		case tea.KeyShiftTab:
			dst := d.directory.Nodes[d.current_node_id].(Node).FocusIndex() - 1
			cmds = append(cmds, d.directory.Nodes[d.current_node_id].(Node).OnFocus(dst))
		case tea.KeyTab:
			dst := d.directory.Nodes[d.current_node_id].(Node).FocusIndex() + 1
			cmds = append(cmds, d.directory.Nodes[d.current_node_id].(Node).OnFocus(dst))
		}
	case base.RegisterMessage:
		if _, ok := msg.Node.(Node); ok {
			_, c = d.directory.Update(msg)
			cmds = append(cmds, c)

			d.dirty = true
		}
	case types.FocusMessage:
		if n, ok := d.directory.Nodes[d.current_node_id]; ok {
			cmds = append(cmds, n.(Node).OnBlur())
		}
		if n, ok := d.directory.Nodes[msg.ID]; ok {
			cmds = append(cmds, n.(Node).OnFocus(msg.Index))
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
			focus_index = d.directory.Nodes[target_id].(Node).NElements() - 1
		} else {
			focus_index = 0
		}
		cmds = append(
			cmds,
			d.directory.Nodes[d.current_node_id].(Node).OnBlur(),
			d.directory.Nodes[target_id].(Node).OnFocus(focus_index),
		)
		d.current_node_id = target_id
	}
	return d, tea.Batch(cmds...)
}

func (d *D) View() string { return d.directory.View() }

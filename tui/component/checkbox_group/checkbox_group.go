package checkbox_group

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/minkezhang/truffle/tui/component/checkbox"
	"github.com/minkezhang/truffle/tui/component/directory/base"
	"github.com/minkezhang/truffle/tui/component/directory/focusable/types"
	"github.com/minkezhang/truffle/tui/component/focusable"
	"github.com/minkezhang/truffle/tui/util/color_profile"
	"github.com/minkezhang/truffle/tui/util/form"
)

type Node struct {
	*focusable.Node

	is_radio bool
	values   []*checkbox.Node
	key      form.Key
}

type O struct {
	Prefix   string
	ParentID string
	IsRadio  bool
	Value    form.Value[[]form.Value[bool]]
}

func New(o O) *Node {
	n := &Node{
		Node:     focusable.New(o.Prefix, o.ParentID, 0),
		is_radio: o.IsRadio,
		values:   []*checkbox.Node{},
		key:      o.Value.Key,
	}
	for _, v := range o.Value.Value {
		b := checkbox.New(checkbox.O{
			Prefix:   fmt.Sprintf("%s-%s", o.Prefix, v.Key.Key),
			ParentID: n.ID(),
			Value:    v,
			IsRadio:  o.IsRadio,
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

func (n *Node) Value() form.Value[[]form.Value[bool]] {
	var values []form.Value[bool]
	for _, v := range n.values {
		values = append(values, v.Value())
	}
	return form.Value[[]form.Value[bool]]{
		Key:   n.key,
		Value: values,
	}
}

func (n *Node) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case checkbox.SelectCheckboxInputMessage:
		// Radio buttons are exclusively selected.
		if msg.ParentID == n.ID() && n.is_radio && msg.Value.Value {
			for _, b := range n.values {
				if b.ID() != msg.ID {
					cmds = append(cmds, func() tea.Msg {
						return checkbox.SelectCheckboxInputMessage{
							ID:       b.ID(),
							ParentID: n.ID(),
							Value: form.Value[bool]{
								Key:   b.Value().Key,
								Value: false,
							},
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

func (n *Node) FocusState() types.FocusState {
	for _, b := range n.values {
		if b.FocusState() == types.FocusStateActive {
			return types.FocusStateActive
		}
	}
	return types.FocusStateNone
}

func (n *Node) View() string {
	var parts []string
	for _, b := range n.values {
		parts = append(parts, b.View())
	}
	return lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.NewStyle().Foreground(
			color_profile.UIForeground[n.FocusState()],
		).Render(n.key.Label),
		strings.Join(parts, " "), // TODO(minkezhang): Render within bounds.
	)
}

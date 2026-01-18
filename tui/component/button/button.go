package button

import (
	"fmt"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lrstanley/bubblezone"
	"github.com/minkezhang/truffle/tui/component/clickable"
	"github.com/minkezhang/truffle/tui/component/directory/base"
	"github.com/minkezhang/truffle/tui/component/directory/focusable/types"
	"github.com/minkezhang/truffle/tui/component/errors"
	"github.com/minkezhang/truffle/tui/component/focusable"
	"github.com/minkezhang/truffle/tui/util/color_profile"
	"github.com/minkezhang/truffle/tui/util/form"
)

type Node struct {
	*focusable.Node

	key       form.Key
	clickable *clickable.Node
}

type O struct {
	ParentID string
	Key      form.Key
}

func New(o O) *Node {
	n := &Node{
		Node: focusable.New("button", o.ParentID, 1),
		key:  o.Key,
	}
	n.clickable = clickable.New(n.ID())
	return n
}

type SubmitMessage struct {
	ID  string
	Key form.Key
}

func (n *Node) Init() tea.Cmd {
	return tea.Sequence(
		n.clickable.Init(),
		func() tea.Msg {
			return base.RegisterMessage{
				Node: n,
			}
		},
	)
}

func (n *Node) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var c tea.Cmd

	if n.FocusState() == types.FocusStateActive {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			if msg.Type == tea.KeyEnter || msg.Type == tea.KeySpace {
				cmds = append(cmds, tea.Sequence(
					func() tea.Msg {
						return SubmitMessage{
							ID:  n.ID(),
							Key: n.key,
						}
					},
					func() tea.Msg {
						return errors.ToLogMessage(
							errors.LevelDebug,
							fmt.Sprintf("%v: clicked button %v", n.ID(), n.key.Key),
						)
					},
				))
			}
		}
	}

	_, c = n.clickable.Update(msg)
	cmds = append(cmds, c)

	switch msg := msg.(type) {
	case clickable.Click:
		if msg.ID == n.clickable.ID() {
			var _cmds []tea.Cmd
			if n.FocusState() == types.FocusStateNone {
				_cmds = append(_cmds, func() tea.Msg {
					return types.FocusMessage{
						ID:    n.ID(),
						Index: 0,
					}
				})
			}
			_cmds = append(_cmds,
				func() tea.Msg {
					return SubmitMessage{
						ID:  n.ID(),
						Key: n.key,
					}
				},
				func() tea.Msg {
					return errors.ToLogMessage(
						errors.LevelDebug,
						fmt.Sprintf("%v: clicked button %v", n.ID(), n.key.Key),
					)
				},
			)
			cmds = append(cmds, tea.Sequence(_cmds...))
		}
	}

	return n, tea.Batch(cmds...)
}

func (n *Node) View() string {
	return zone.Mark(
		n.clickable.ID(),
		lipgloss.NewStyle().Border(lipgloss.NormalBorder()).Foreground(
			color_profile.UIForeground[n.FocusState()]).BorderForeground(
			color_profile.UIForeground[n.FocusState()]).Padding(0, 1, 0, 1).Render(
			n.key.Label,
		),
	)

}

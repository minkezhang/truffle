package checkbox

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

	is_radio      bool
	clickable     *clickable.Node
	clickable_box *clickable.Node
	value         form.Value[bool]
}

type O struct {
	Prefix   string
	ParentID string
	Value    form.Value[bool]
	IsRadio  bool
}

type SelectCheckboxInputMessage struct {
	ID       string
	ParentID string
	Value    form.Value[bool]
}

func New(o O) *Node {
	n := &Node{
		Node:     focusable.New(o.Prefix, o.ParentID, 1),
		value:    o.Value,
		is_radio: o.IsRadio,
	}
	n.clickable = clickable.New(n.ID())
	n.clickable_box = clickable.New(n.ID())
	return n
}

func (n *Node) Value() form.Value[bool] { return n.value }

func (n *Node) Init() tea.Cmd {
	return tea.Sequence(
		tea.Batch(
			n.clickable.Init(),
			n.clickable_box.Init(),
		),
		func() tea.Msg {
			return base.RegisterMessage{
				Node: n,
			}
		},
	)
}

func (n *Node) toggle() tea.Cmd {
	return tea.Batch(
		func() tea.Msg {
			return SelectCheckboxInputMessage{
				ID:       n.ID(),
				ParentID: n.ParentID(),
				Value: form.Value[bool]{
					Key:   n.value.Key,
					Value: !n.value.Value || n.is_radio, // can't manually deselect radio button
				},
			}
		},
		func() tea.Msg {
			return errors.ToLogMessage(
				errors.LevelDebug,
				fmt.Sprintf("%v: selecting value %v = %v", n.ID(), n.value.Key.Key, !n.value.Value || n.is_radio),
			)
		},
	)
}

func (n *Node) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var c tea.Cmd

	switch msg := msg.(type) {
	case SelectCheckboxInputMessage:
		if msg.ID == n.ID() {
			n.value = msg.Value
		}
	}
	if n.FocusState() == types.FocusStateActive {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			if msg.Type == tea.KeyEnter || msg.Type == tea.KeySpace {
				cmds = append(cmds, n.toggle())
			}
		case clickable.Click:
			if msg.ID == n.clickable_box.ID() {
				cmds = append(cmds, n.toggle())
			}
		}
	} else {
		switch msg := msg.(type) {
		case clickable.Click:
			if msg.ID == n.clickable.ID() {
				cmds = append(cmds, func() tea.Msg {
					return types.FocusMessage{
						ID:    n.ID(),
						Index: 0,
					}
				})
			}
			if msg.ID == n.clickable_box.ID() {
				cmds = append(cmds, n.toggle())
			}
		}
		_, c = n.clickable.Update(msg)
		cmds = append(cmds, c)
	}
	_, c = n.clickable_box.Update(msg)
	cmds = append(cmds, c)

	return n, tea.Batch(cmds...)
}

func (n *Node) View() string {
	return zone.Mark(
		n.clickable.ID(),
		lipgloss.NewStyle().Foreground(
			color_profile.UIForeground[n.FocusState()],
		).Render(
			fmt.Sprintf(
				"%s %s",
				zone.Mark(
					n.clickable_box.ID(),
					fmt.Sprintf("%s%s%s",
						map[bool]string{
							true:  "(",
							false: "[",
						}[n.is_radio],
						map[bool]string{
							true: map[bool]string{
								true:  "●",
								false: "■",
							}[n.is_radio],
							false: " ",
						}[n.value.Value],
						map[bool]string{
							true:  ")",
							false: "]",
						}[n.is_radio],
					),
				),
				n.value.Key.Label,
			),
		),
	)
}

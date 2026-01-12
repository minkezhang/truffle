package checkbox

import (
	"fmt"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lrstanley/bubblezone"
	"github.com/minkezhang/truffle/tui/component/clickable"
	"github.com/minkezhang/truffle/tui/component/directory/base"
	"github.com/minkezhang/truffle/tui/component/directory/focusable/types"
	"github.com/minkezhang/truffle/tui/component/focusable"
	"github.com/minkezhang/truffle/tui/util/color_profile"
)

type Node struct {
	*focusable.Node

	is_radio      bool
	label         string
	value         string
	is_selected   bool
	clickable     *clickable.Node
	clickable_box *clickable.Node
}

type O struct {
	Prefix     string
	ParentID   string
	Label      string
	Value      string
	IsSelected bool
	IsRadio    bool
}

type SelectCheckboxInput struct {
	ID         string
	ParentID   string
	Value      string
	IsSelected bool
}

func New(o O) *Node {
	n := &Node{
		Node:        focusable.New(o.Prefix, o.ParentID, 1),
		label:       o.Label,
		value:       o.Value,
		is_selected: o.IsSelected,
		is_radio:    o.IsRadio,
	}
	n.clickable = clickable.New(n.ID())
	n.clickable_box = clickable.New(n.ID())
	return n
}

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

func (n *Node) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var c tea.Cmd

	switch msg := msg.(type) {
	case SelectCheckboxInput:
		if msg.ID == n.ID() {
			n.is_selected = msg.IsSelected
		}
	}
	if n.FocusState() == types.FocusStateActive {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			if msg.Type == tea.KeyEnter || msg.Type == tea.KeySpace {
				cmds = append(cmds, func() tea.Msg {
					return SelectCheckboxInput{
						ID:         n.ID(),
						ParentID:   n.ParentID(),
						IsSelected: !n.is_selected || n.is_radio,  // can't manually deselect radio
						Value:      n.value,
					}
				})
			}
		case clickable.Click:
			if msg.ID == n.clickable_box.ID() {
				cmds = append(cmds, func() tea.Msg {
					return SelectCheckboxInput{
						ID:         n.ID(),
						ParentID:   n.ParentID(),
						IsSelected: !n.is_selected || n.is_radio,  // can't manually deselect radio
						Value:      n.value,
					}
				})
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
				cmds = append(cmds, func() tea.Msg {
					return SelectCheckboxInput{
						ID:         n.ID(),
						ParentID:   n.ParentID(),
						IsSelected: !n.is_selected || n.is_radio,  // can't manually deselect radio
						Value:      n.value,
					}
				})
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
								true:  "o",
								false: "x",
							}[n.is_radio],
							false: " ",
						}[n.is_selected],
						map[bool]string{
							true:  ")",
							false: "]",
						}[n.is_radio],
					),
				),
				n.label,
			),
		),
	)
}

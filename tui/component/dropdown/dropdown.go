package dropdown

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"github.com/lrstanley/bubblezone"
	"github.com/minkezhang/truffle/tui/component/choices"
	"github.com/minkezhang/truffle/tui/component/clickable"
	"github.com/minkezhang/truffle/tui/component/column"
	"github.com/minkezhang/truffle/tui/component/directory/base"
	"github.com/minkezhang/truffle/tui/component/directory/focusable/types"
	"github.com/minkezhang/truffle/tui/component/focusable"
	"github.com/minkezhang/truffle/tui/util/color_profile"
	"github.com/minkezhang/truffle/tui/util/form"
)

type O struct {
	Prefix   string
	ParentID string
	Column   *column.C
	Width    int
	Key      form.Key
	Choices  []form.Key
	Prompt   string
	Height   int
}

type Node struct {
	*focusable.Node

	column    *column.C
	prompt    string
	clickable *clickable.Node
	choices   *choices.Node
}

func New(o O) *Node {
	n := &Node{
		Node:   focusable.New(o.Prefix, o.ParentID, 1),
		column: o.Column,
		prompt: o.Prompt,
	}
	n.clickable = clickable.New(n.ID())
	n.choices = choices.New(choices.O{
		ParentID: n.ID(),
		Width:    o.Column.Content(),
		Height:   o.Height,
		Key:      o.Key,
		Choices:  o.Choices,
	})
	return n
}

func (n *Node) Init() tea.Cmd {
	return tea.Sequence(
		n.clickable.Init(),
		n.choices.Init(),
		n.choices.SetIsInvisible(true),
		func() tea.Msg {
			return base.RegisterMessage{
				Node: n,
			}
		},
	)
}

func (n *Node) SetValue(vs []form.Key, index int) tea.Cmd { return n.choices.SetValue(vs, index) }
func (n *Node) Value() form.Value[form.Key]               { return n.choices.Value() }

func (n *Node) do_expand() tea.Cmd {
	is_expanded := !n.choices.IsInvisible()

	cmds := []tea.Cmd{
		n.choices.SetIsInvisible(!n.choices.IsInvisible()),
	}

	if !is_expanded {
		cmds = append(cmds,
			func() tea.Msg {
				return types.FocusMessage{
					ID: n.choices.ID(),
				}
			},
		)
	}
	return tea.Sequence(cmds...)
}

func (n *Node) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if n.FocusState() == types.FocusStateActive {
			if msg.Type == tea.KeyEnter || msg.Type == tea.KeySpace {
				cmds = append(cmds, n.do_expand())
			}
		}
	case clickable.Click:
		if msg.ID == n.clickable.ID() {
			cmds = append(cmds, n.do_expand())
		}
	}

	for _, m := range []tea.Model{
		n.clickable,
		n.choices,
	} {
		_, c := m.Update(msg)
		cmds = append(cmds, c)
	}

	return n, tea.Batch(cmds...)
}

func (n *Node) View() string {
	var choices string
	if !n.choices.IsInvisible() {
		choices = n.choices.View()
	}
	c := n.column.WithWidth(n.column.Width() - /* is_expanded */ 3)
	return n.column.RenderOrDie(
		lipgloss.JoinVertical(
			lipgloss.Left,
			zone.Mark(
				n.clickable.ID(),
				lipgloss.JoinHorizontal(
					lipgloss.Top,
					c.Style().Render(
						lipgloss.JoinVertical(
							lipgloss.Left,
							c.Style().MaxWidth(c.Width()).Inline(true).Render(
								fmt.Sprintf(
									"%v%v",
									lipgloss.NewStyle().Foreground(
										color_profile.UIForeground[n.FocusState()],
									).Render(n.prompt),
									ansi.Truncate(
										n.Value().Value.Label,
										c.Content()-lipgloss.Width(n.prompt),
										"…",
									),
								),
							),
							lipgloss.NewStyle().Foreground(
								color_profile.UIForeground[n.FocusState()],
							).Render(
								fmt.Sprintf(
									"%v%v%v%v",
									map[bool]string{
										true:  "──",
										false: "─ ",
									}[n.Value().Key.Label == ""],
									n.Value().Key.Label,
									map[bool]string{
										true:  "─",
										false: " ",
									}[n.Value().Key.Label == ""],
									strings.Repeat("─", c.Content()-len(n.Value().Key.Label)-3),
								),
							),
						),
					),
					lipgloss.NewStyle().Width(3).Border(lipgloss.NormalBorder(), false, false, true, false).BorderForeground(
						color_profile.UIForeground[n.FocusState()],
					).Foreground(
						color_profile.UIForeground[n.FocusState()],
					).Render(
						map[bool]string{
							false: " ↓ ",
							true:  " → ",
						}[n.choices.IsInvisible()],
					),
				),
			),
			choices,
		),
	)
}

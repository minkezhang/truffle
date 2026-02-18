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

	column     *column.C
	prompt     string
	clickable  *clickable.Node
	expandable *clickable.Node
	choices    *choices.Node
}

func New(o O) *Node {
	n := &Node{
		Node:   focusable.New(o.Prefix, o.ParentID, 2),
		column: o.Column,
		prompt: o.Prompt,
	}
	n.clickable = clickable.New(n.ID())
	n.expandable = clickable.New(n.ID())
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
		n.expandable.Init(),
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

	if is_expanded {
		cmds = append(cmds,
			func() tea.Msg {
				return types.FocusMessage{
					ID:    n.ID(),
					Index: 0,
				}
			},
		)
	} else {
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
	case clickable.Click:
		if msg.ID == n.expandable.ID() {
			cmds = append(cmds, n.do_expand())
		}
	case tea.KeyMsg:
		if n.FocusIndex() == 0 && n.FocusState() == types.FocusStateActive {
			if msg.Type == tea.KeyEnter || msg.Type == tea.KeySpace {
				cmds = append(cmds, n.do_expand())
			}
		}

	}

	if n.FocusState() != types.FocusStateActive {
		switch msg := msg.(type) {
		case clickable.Click:
			if msg.ID == n.clickable.ID() {
				cmds = append(cmds, func() tea.Msg {
					return types.FocusMessage{
						ID:    n.ID(),
						Index: 1,
					}
				})
			}
		}
	}

	for _, m := range []tea.Model{
		n.clickable,
		n.expandable,
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
	c := n.column.WithWidth(n.column.Width() - 3)
	return n.column.RenderOrDie(
		lipgloss.JoinVertical(
			lipgloss.Left,
			lipgloss.JoinHorizontal(
				lipgloss.Top,
				zone.Mark(
					n.clickable.ID(),
					c.Style().Render(
						lipgloss.JoinVertical(
							lipgloss.Left,
							c.Style().MaxWidth(c.Width()).Inline(true).Render(
								fmt.Sprintf(
									"%v%v",
									lipgloss.NewStyle().Foreground(
										map[bool]lipgloss.Color{
											true:  color_profile.UIForeground[types.FocusStateActive],
											false: color_profile.UIForeground[types.FocusStateNone],
										}[n.FocusState() == types.FocusStateActive && n.FocusIndex() == 1],
									).Render(n.prompt),
									ansi.Truncate(
										n.Value().Value.Label,
										c.Content()-lipgloss.Width(n.prompt),
										"…",
									),
								),
							),
							lipgloss.NewStyle().Foreground(
								map[bool]lipgloss.Color{
									true:  color_profile.UIForeground[types.FocusStateActive],
									false: color_profile.UIForeground[types.FocusStateNone],
								}[n.FocusState() == types.FocusStateActive && n.FocusIndex() == 1],
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
				),
				zone.Mark(
					n.expandable.ID(),
					lipgloss.NewStyle().Width(3).Border(lipgloss.NormalBorder(), false, false, true, false).BorderForeground(
						map[bool]lipgloss.Color{
							true:  color_profile.UIForeground[types.FocusStateActive],
							false: color_profile.UIForeground[types.FocusStateNone],
						}[n.FocusState() == types.FocusStateActive && n.FocusIndex() == 1],
					).Foreground(
						map[bool]lipgloss.Color{
							true:  color_profile.UIForeground[types.FocusStateActive],
							false: color_profile.UIForeground[types.FocusStateNone],
						}[n.FocusState() == types.FocusStateActive && n.FocusIndex() == 0],
					).Render(
						map[bool]string{
							true:  "(+)",
							false: "(-)",
						}[n.choices.IsInvisible()],
					),
				),
			),
			choices,
		),
	)
}

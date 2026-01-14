package textinput

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
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

	max_width int
	input     textinput.Model
	clickable *clickable.Node
	key       form.Key
}

type O struct {
	Prefix      string
	ParentID    string
	Width       int
	Placeholder string
	Prompt      string
	Value       form.Value[string]
}

func New(o O) *Node {
	t := textinput.New()
	t.Width = o.Width - len(o.Prompt) - 1 // cursor
	t.Prompt = o.Prompt
	t.Placeholder = o.Placeholder
	t.SetValue(o.Value.Value)

	n := &Node{
		Node:      focusable.New(o.Prefix, o.ParentID, 1),
		input:     t,
		max_width: o.Width,
		key:       o.Value.Key,
	}
	n.clickable = clickable.New(n.ID())
	return n
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

type SubmitTextInput struct {
	ID    string
	Value form.Value[string]
}

func (n *Node) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var c tea.Cmd

	if n.FocusState() == types.FocusStateActive {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			n.input, c = n.input.Update(msg)
			cmds = append(cmds, c)

			switch msg.Type {
			case tea.KeyEnter:
				v := n.input.Value()
				n.input.SetValue("") // TODO(minkezhang): Only if n.accept_enter
				cmds = append(
					cmds,
					func() tea.Msg {
						return SubmitTextInput{
							ID: n.ID(),
							Value: form.Value[string]{
								Key:   n.key,
								Value: v,
							},
						}
					},
					func() tea.Msg {
						return errors.ToLogMessage(
							errors.LevelDebug,
							fmt.Sprintf("%v: submitting value %v = \"%v\"", n.ID(), n.key.Key, v),
						)
					},
				)
			}
		}
		n.input.Cursor, c = n.input.Cursor.Update(msg)
		cmds = append(cmds, c)
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
		}
		_, c = n.clickable.Update(msg)
		cmds = append(cmds, c)
	}
	return n, tea.Batch(cmds...)
}

func (n *Node) View() string {
	return zone.Mark(
		n.clickable.ID(),
		lipgloss.JoinVertical(
			lipgloss.Left,
			lipgloss.NewStyle().Width(n.max_width).MaxWidth(n.max_width).Render(n.input.View()),
			lipgloss.NewStyle().Foreground(color_profile.UIForeground[n.FocusState()]).Render(
				fmt.Sprintf(
					"%v%v%v",
					n.key.Label,
					map[bool]string{
						true:  "─",
						false: " ",
					}[n.key.Label == ""],
					strings.Repeat("─", n.max_width-len(n.key.Label)-1),
				),
			),
		),
	)
}

func (n *Node) OnFocus(i int) tea.Cmd {
	cmds := []tea.Cmd{n.Node.OnFocus(i)}
	if n.FocusState() == types.FocusStateActive {
		n.input.PromptStyle = n.input.PromptStyle.Foreground(
			color_profile.UIForeground[types.FocusStateActive],
		)
		cmds = append(cmds,
			n.input.Focus(),
			n.input.Cursor.Focus(),
		)
	}
	return tea.Sequence(cmds...)
}

func (n *Node) OnBlur() tea.Cmd {
	n.input.Blur()
	n.input.Cursor.Blur()
	n.input.PromptStyle = n.input.PromptStyle.Foreground(
		color_profile.UIForeground[types.FocusStateNone],
	)
	return tea.Sequence(
		n.Node.OnBlur(),
	)
}

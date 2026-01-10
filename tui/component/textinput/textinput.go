package textinput

import (
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbletea"
	"github.com/lrstanley/bubblezone"
	"github.com/minkezhang/truffle/tui/component/clickable"
	"github.com/minkezhang/truffle/tui/component/focusable/directory/types"
	"github.com/minkezhang/truffle/tui/component/focusable"
)

type Node struct {
	*focusable.Node

	input     textinput.Model
	clickable *clickable.Node
}

type O struct {
	Prefix      string
	ParentID    string
	Width       int
	Placeholder string
	Prompt      string
}

func New(o O) *Node {
	t := textinput.New()
	t.Width = t.Width
	t.Prompt = t.Prompt
	t.Placeholder = t.Placeholder

	n := &Node{
		Node:  focusable.New(o.Prefix, o.ParentID, 1),
		input: t,
	}
	n.clickable = clickable.New(n.ID())
	return n
}

func (n *Node) Init() tea.Cmd {
	return tea.Sequence(
		n.clickable.Init(),
		func() tea.Msg {
			return types.RegisterNodeMessage{
				Node: n,
			}
		},
	)
}

type SubmitTextInput struct {
	ID    string
	Value string
}

func (n *Node) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var c tea.Cmd

	if n.FocusState() == types.FocusStateActive {
		switch msg := msg.(type) {
		case clickable.Click:
			if msg.ID == n.clickable.ID() {
				n.input.SetValue("")
			}
		case tea.KeyMsg:
			n.input, c = n.input.Update(msg)
			cmds = append(cmds, c)

			switch msg.Type {
			case tea.KeyEnter:
				v := n.input.Value()
				n.input.SetValue("")
				cmds = append(
					cmds,
					func() tea.Msg {
						return SubmitTextInput{
							ID:    n.ID(),
							Value: v,
						}
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
			_, c := n.clickable.Update(msg)
			cmds = append(cmds, c)
		}
	}
	return n, tea.Batch(cmds...)
}

func (n *Node) View() string {
	return zone.Mark(n.clickable.ID(), n.input.View())
}

func (n *Node) OnFocus(i int) tea.Cmd {
	cmds := []tea.Cmd{n.Node.OnFocus(i)}
	if n.FocusState() == types.FocusStateActive {
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
	return tea.Sequence(
		n.Node.OnBlur(),
	)
}

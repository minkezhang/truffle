package textinput

import (
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbletea"
	"github.com/minkezhang/truffle/tui/components/directory/types"
	"github.com/minkezhang/truffle/tui/components/focusable"
)

type Node struct {
	focusable.Node

	input textinput.Model
}

func New(prefix string, parent_id string) *Node {
	return &Node{
		Node: focusable.New(prefix, parent_id, 1),
	}
}

func (n *Node) Init() tea.Cmd {
	return func() tea.Msg {
		return types.RegisterNodeMessage{
			Node: n,
		}
	}
}

type SubmitTextInput struct {
	ID string
	Value string
}

func (n *Node) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	if n.FocusState() == types.FocusStateActive {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			// Mouse input is being passed into search input as KeyMsg; it
			// is unclear how or why this is happening.
			//
			// Such a KeyMsg is of the form
			//
			//   { Type: KeyRunes, Alt: true, Runes: []rune{'['} }
			//
			// or
			//
			//   { Type: KeyRunes, Alt: false, Runes: []rune{...} }
			if (len(msg.Runes) <= 1 && !msg.Alt) || msg.Paste {
				var c tea.Cmd
				n.input, c = n.input.Update(msg)
				cmds = append(cmds, c)
			}

			switch msg.Type {
			case tea.KeyEnter:
				v := n.input.Value()
				n.input.SetValue("")
				cmds = append(
					cmds,
					func() tea.Msg {
						return SubmitTextInput{
							ID: n.ID(),
							Value:  v,
						}
					},
				)
			}
		}
	}
	return n, nil
}

func (n *Node) View() string { return n.input.View() }

func (n *Node) OnFocus(i int) tea.Cmd {
	cmds := []tea.Cmd{n.Node.OnFocus(i)}
	if n.FocusState() == types.FocusStateActive {
		cmds = append(cmds, n.input.Focus())
	}
	return tea.Sequence(cmds...)
}
func (n *Node) OnBlur() tea.Cmd {
	n.input.Blur()
	return n.Node.OnBlur()
}

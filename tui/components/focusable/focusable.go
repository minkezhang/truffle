package focusable

import (
	"fmt"

	"github.com/charmbracelet/bubbletea"
	"github.com/lrstanley/bubblezone"
	"github.com/minkezhang/truffle/tui/components/directory/types"
)

type Node struct {
	id          string
	parent_id   string
	n_elements  int
	index       int
	focus_state types.FocusState
}

func New(prefix string, parent_id string, n_elements int) Node {
	if prefix == "" {
		prefix = "none"
	}
	return Node{
		id:         fmt.Sprintf("%v:%v", prefix, zone.NewPrefix()),
		parent_id:  parent_id,
		n_elements: n_elements,
	}
}

func (n *Node) ID() string                   { return n.id }
func (n *Node) ParentID() string             { return n.parent_id }
func (n *Node) FocusState() types.FocusState { return n.focus_state }
func (n *Node) FocusIndex() int              { return n.index }
func (n *Node) NElements() int               { return n.n_elements }

func (n *Node) OnFocus(i int) tea.Cmd {
	n.focus_state = types.FocusStateActive
	if i < 0 || i >= n.NElements() {
		return func() tea.Msg {
			return types.EOFMessage{
				ID:     n.ID(),
				IsHead: i < 0,
			}
		}
	}
	n.index = i
	return nil
}

func (n *Node) OnBlur() tea.Cmd {
	n.focus_state = types.FocusStateNone
	return nil
}

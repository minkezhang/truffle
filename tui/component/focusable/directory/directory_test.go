package directory

import (
	"os"
	"testing"

	"github.com/charmbracelet/bubbletea"
	"github.com/google/go-cmp/cmp"
	"github.com/lrstanley/bubblezone"
	"github.com/minkezhang/truffle/tui/component/directory/types"
	"github.com/minkezhang/truffle/tui/component/focusable"
)

var _ types.Node = &Mock{}

type Mock struct {
	*focusable.Node
}

func (m *Mock) Init() tea.Cmd {
	return func() tea.Msg {
		return types.RegisterNodeMessage{
			Node: m,
		}
	}
}

func (m *Mock) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return m, nil }
func (m *Mock) View() string                            { return "" }
func (m *Mock) OnFocus(i int) tea.Cmd                   { return m.Node.OnFocus(i) } // May be overridden
func (m *Mock) OnBlur() tea.Cmd                         { return m.Node.OnBlur() }   // May be overridden

func TestMain(m *testing.M) {
	zone.NewGlobal()
	os.Exit(m.Run())
}

// Consider the following node tree --
//
//	  A
//	 / \
//	B   C
//	|
//	D
//
// order() returns a pre-order ordering of the nodes; we expect therefore a
// result of [A, B, D, C]
func TestOrder(t *testing.T) {
	na := &Mock{Node: focusable.New("test-a", "", 0)}
	nb := &Mock{Node: focusable.New("test-b", na.ID(), 0)}
	nc := &Mock{Node: focusable.New("test-c", na.ID(), 0)}
	nd := &Mock{Node: focusable.New("test-d", nb.ID(), 0)}
	d := &D{
		nodes: map[string]types.Node{
			na.ID(): na,
			nb.ID(): nb,
			nc.ID(): nc,
			nd.ID(): nd,
		},
		children: map[string][]string{
			"":      []string{na.ID()},
			na.ID(): []string{nb.ID(), nc.ID()},
			nb.ID(): []string{nd.ID()},
			nc.ID(): []string{},
			nd.ID(): []string{},
		},
		parent: map[string]string{
			na.ID(): "",
			nb.ID(): nb.ParentID(),
			nc.ID(): nc.ParentID(),
			nd.ID(): nd.ParentID(),
		},
		dirty: true,
	}

	got, want := d.order(), []string{na.ID(), nb.ID(), nd.ID(), nc.ID()}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("order() mismatch (-want +got):\n%v", diff)
	}
}

func TestUpdate(t *testing.T) {
	t.Run("RegisterNodeMessage", func(t *testing.T) {
		na := &Mock{Node: focusable.New("test-a", "", 0)}
		nb := &Mock{Node: focusable.New("test-b", na.ID(), 0)}
		nc := &Mock{Node: focusable.New("test-c", na.ID(), 0)}
		nd := &Mock{Node: focusable.New("test-d", nb.ID(), 0)}
		d := New()
		d.Update(types.RegisterNodeMessage{Node: na})
		d.Update(types.RegisterNodeMessage{Node: nb})
		d.Update(types.RegisterNodeMessage{Node: nc})
		d.Update(types.RegisterNodeMessage{Node: nd})
		got, want := d.order(), []string{na.ID(), nb.ID(), nd.ID(), nc.ID()}
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("order() mismatch (-want +got):\n%v", diff)
		}
	})
	t.Run("KeyMsg", func(t *testing.T) {
		t.Run("Simple/KeyTab", func(t *testing.T) {
			n := &Mock{Node: focusable.New("test", "", 2)}
			d := New()
			d.Update(types.RegisterNodeMessage{Node: n})
			d.Update(types.FocusMessage{ID: n.ID(), Index: 0}) // Manually simulate current focus
			d.Update(tea.KeyMsg{Type: tea.KeyTab})
			if got, want := n.FocusIndex(), 1; got != want {
				t.Errorf("FocusIndex() = %v, want = %v", got, want)
			}
		})
		t.Run("Simple/KeyTab/EOF", func(t *testing.T) {
			n := &Mock{Node: focusable.New("test", "", 2)}
			d := New()
			d.Update(types.RegisterNodeMessage{Node: n})
			d.Update(types.FocusMessage{ID: n.ID(), Index: 0}) // Manually simulate current focus
			d.Update(tea.KeyMsg{Type: tea.KeyTab})
			d.Update(tea.KeyMsg{Type: tea.KeyTab}) // EOF
			if got, want := n.FocusIndex(), 1; got != want {
				t.Errorf("FocusIndex() = %v, want = %v", got, want)
			}
		})
	})
	t.Run("FocusMessage", func(t *testing.T) {

	})
	t.Run("EOFMessage", func(t *testing.T) {
		t.Run("Simple/Next", func(t *testing.T) {
			na := &Mock{Node: focusable.New("test-a", "", 1)}
			nb := &Mock{Node: focusable.New("test-b", "", 1)}
			d := New()
			d.Update(types.RegisterNodeMessage{Node: na})
			d.Update(types.RegisterNodeMessage{Node: nb})
			d.Update(types.FocusMessage{ID: na.ID(), Index: 0})
			d.Update(types.EOFMessage{ID: na.ID(), IsHead: false})
			if got, want := d.current_node_id, nb.ID(); got != want {
				t.Errorf("current_node_id = %v, want = %v", got, want)
			}
		})
		t.Run("Simple/Next/Cycle", func(t *testing.T) {
			na := &Mock{Node: focusable.New("test-a", "", 1)}
			nb := &Mock{Node: focusable.New("test-b", "", 1)}
			d := New()
			d.Update(types.RegisterNodeMessage{Node: na})
			d.Update(types.RegisterNodeMessage{Node: nb})
			d.Update(types.FocusMessage{ID: nb.ID(), Index: 0})
			d.Update(types.EOFMessage{ID: nb.ID(), IsHead: false})
			if got, want := d.current_node_id, na.ID(); got != want {
				t.Errorf("current_node_id = %v, want = %v", got, want)
			}
		})
	})
}

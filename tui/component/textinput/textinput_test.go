package textinput

import (
	"os"
	"testing"

	"github.com/charmbracelet/bubbletea"
	"github.com/lrstanley/bubblezone"
	"github.com/minkezhang/truffle/tui/component/directory/focusable/types"
)

func TestMain(m *testing.M) {
	zone.NewGlobal()
	os.Exit(m.Run())
}

func TestInput(t *testing.T) {
	t.Run("NotFocused", func(t *testing.T) {
		n := New(O{})
		n.Init()()
		if got, want := n.FocusState(), types.FocusStateNone; got != want {
			t.Errorf("FocusState() = %v, want = %v", got, want)
		}
		_, c := n.Update(tea.KeyMsg{
			Type:  tea.KeyRunes,
			Runes: []rune{'f'},
		})
		if c != nil {
			c()
		}
		if got, want := n.input.Value(), ""; got != want {
			t.Errorf("Value() = %v, want = %v", got, want)
		}
	})
	t.Run("Focused", func(t *testing.T) {
		n := New(O{})
		n.Init()()
		n.OnFocus(0)()
		if got, want := n.FocusState(), types.FocusStateActive; got != want {
			t.Errorf("FocusState() = %v, want = %v", got, want)
		}
		_, c := n.Update(tea.KeyMsg{
			Type:  tea.KeyRunes,
			Runes: []rune{'f'},
		})
		if c != nil {
			c()
		}
		if got, want := n.input.Value(), "f"; got != want {
			t.Errorf("Value() = %v, want = %v", got, want)
		}
	})
}

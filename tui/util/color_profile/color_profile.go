package color_profile

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/minkezhang/truffle/tui/component/directory/focusable/types"
)

var (
	UIForeground = map[types.FocusState]lipgloss.Color{
		types.FocusStateNone:   lipgloss.Color("8"),
		types.FocusStateActive: lipgloss.Color("7"),
	}

	SupplementaryText = lipgloss.Color("8")
	UserViewText      = lipgloss.Color("5")
)

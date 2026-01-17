package color_profile

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/minkezhang/truffle/tui/component/directory/focusable/types"
	"github.com/minkezhang/truffle/tui/component/errors"
)

var (
	UIForeground = map[types.FocusState]lipgloss.Color{
		types.FocusStateNone:   lipgloss.Color("8"),
		types.FocusStateActive: lipgloss.Color("7"),
	}

	ForegroundInverted   = lipgloss.Color("0")
	ForegroundNegligible = lipgloss.Color("8")
	ForegroundNormal     = lipgloss.Color("7")
	ForegroundImportant  = lipgloss.Color("6")
	ForegroundCritical   = lipgloss.Color("1")

	BackgroundNegligible = lipgloss.Color("8")

	SupplementaryText = lipgloss.Color("8")
	SupplementaryUI   = lipgloss.Color("8")
	UserViewText      = lipgloss.Color("5")

	LogForeground = map[errors.Level]lipgloss.Color{
		errors.LevelDebug: lipgloss.Color("8"),
		errors.LevelInfo:  lipgloss.Color("7"),
		errors.LevelWarn:  lipgloss.Color("3"),
		errors.LevelError: lipgloss.Color("1"),
	}
)

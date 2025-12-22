package input

import (
	"github.com/charmbracelet/bubbletea"
)

func IsMouseJustPressed(msg tea.Msg) bool {
	switch msg := msg.(type) {
	case tea.MouseMsg:
		return msg.Button == tea.MouseButtonLeft && msg.Action == tea.MouseActionPress
	default:
		return false
	}
}

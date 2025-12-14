package book

import (
	"github.com/lrstanley/bubblezone"
	"github.com/charmbracelet/bubbletea"
)

func (m *M) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.MouseMsg:
		m.hover = zone.Get("TEST").InBounds(msg)
	}
	return m, nil
}

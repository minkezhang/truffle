package book

import (
	"github.com/charmbracelet/bubbletea"
	"github.com/lrstanley/bubblezone"
)

func (m *M) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.MouseMsg:
		m.hover = zone.Get("TEST").InBounds(msg)
	}
	return m, nil
}

package root

import (
	"github.com/charmbracelet/bubbletea"
)

func (m *M) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.debug.Update(msg)

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			fallthrough
		case tea.KeyCtrlD:
			return m, tea.Quit
		case tea.KeyCtrlZ:
			return m, tea.Suspend
		case tea.KeyTab:
			m.overlay = !m.overlay
			return m, nil
		}
	}
	return m, nil
}

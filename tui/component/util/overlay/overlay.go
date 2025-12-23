package overlay

import (
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/yorukot/superfile/src/pkg/string_function"
)

type M struct {
	FG tea.Model
	BG tea.Model
}

func (m M) Init() tea.Cmd { return nil }

func (m M) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return m, nil }

func (m M) View() string {
	fg := m.FG.View()
	bg := m.BG.View()
	w := lipgloss.Width(fg)
	h := lipgloss.Height(fg)
	if x := lipgloss.Width(bg); x > w {
		w = x
	}
	if i := lipgloss.Height(bg); i > h {
		h = i
	}
	return stringfunction.PlaceOverlay(
		0,
		0,
		lipgloss.Place(lipgloss.Width(fg), h, lipgloss.Top, lipgloss.Left, fg),
		lipgloss.Place(lipgloss.Width(bg), h, lipgloss.Top, lipgloss.Left, bg),
	)
}

func Overlay(fg string, bg string) string {
	w := lipgloss.Width(fg)
	h := lipgloss.Height(fg)
	if x := lipgloss.Width(bg); x > w {
		w = x
	}
	if i := lipgloss.Height(bg); i > h {
		h = i
	}
	return stringfunction.PlaceOverlay(
		0,
		0,
		lipgloss.Place(lipgloss.Width(fg), lipgloss.Height(fg), lipgloss.Top, lipgloss.Left, fg),
		lipgloss.Place(lipgloss.Width(bg), lipgloss.Height(bg), lipgloss.Top, lipgloss.Left, bg),
	)
}

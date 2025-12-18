package search

import (
	"fmt"
	"log/slog"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type O struct{}

type M struct {
	search textinput.Model
}

func New(o O) *M {
	m := &M{
		search: textinput.New(),
	}
	m.search.Width = 30
	m.search.Placeholder = "The Apothecary Diaries"
	m.search.Prompt = "⚲ "
	m.search.SetSuggestions([]string{"Apothecary Diaries", "ABBA", "Escaflowne"})
	return m
}

func (m *M) Init() tea.Cmd { return m.search.Focus() }

type QueryMsg string

func (m *M) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := []tea.Cmd{}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			slog.Debug(fmt.Sprintf("sending query request to root: %v\n", m.search.Value()))
			cmds = append(cmds, func() tea.Msg { return QueryMsg(m.search.Value()) })
		}
	}

	var c tea.Cmd
	m.search, c = m.search.Update(msg)

	cmds = append(cmds, c)

	return m, tea.Batch(cmds...)

}

func (m *M) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.NewStyle().Border(lipgloss.NormalBorder(), false, false, true, false).Padding(0, 1).Render(
			m.search.View(),
		),
	)
}

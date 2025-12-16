package tab

import (
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lrstanley/bubblezone"
)

type State int

const (
	StateNone State = iota
	StateActive
)

type O struct {
	Value string
	State State
}

type M struct {
	key   string
	value string
	state State
}

func New(o O) *M {
	return &M{
		key:   zone.NewPrefix(),
		value: o.Value,
		state: o.State,
	}
}

type FocusMsg struct {
	Key    string
	Target State
}

func (m *M) Key() string { return m.key }

func (m *M) Init() tea.Cmd { return nil }

func (m *M) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case FocusMsg:
		if msg.Key == m.key {
			m.state = msg.Target
		}
	case tea.MouseMsg:
		if zone.Get(m.key).InBounds(msg) && msg.Button == tea.MouseButtonLeft {
			m.state = StateActive
		}
		if msg.Button == tea.MouseButtonRight {
			m.state = StateNone
		}
	}

	return m, nil
}

var (
	borders = map[State]lipgloss.Style{
		StateNone: lipgloss.NewStyle().Padding(0, 1).Margin(0, 1).Border(
			lipgloss.NormalBorder(), false, false, true, false,
		).BorderForeground(lipgloss.Color("8")).Foreground(lipgloss.Color("8")),
		StateActive: lipgloss.NewStyle().Padding(0, 1).Margin(0, 1).Border(
			lipgloss.NormalBorder(), false, false, true, false,
		).BorderForeground(lipgloss.Color("6")).Foreground(lipgloss.Color("6")),
	}
)

func (m *M) View() string {
	return zone.Mark(m.key, borders[m.state].Render(m.value))
}

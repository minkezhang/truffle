package button

import (
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lrstanley/bubblezone"

	model_ui "github.com/minkezhang/truffle/tui/component/util/model"
)

type State int

const (
	StateInactive State = iota
	StateActive
)

type O struct {
	model_ui.O

	Name string
}

type M struct {
	*model_ui.Base

	name  string
	state State
}

func New(o O) *M {
	return &M{
		Base: model_ui.New(o.O),
		name: o.Name,
	}
}

func (m *M) Init() tea.Cmd { return nil }

func (m *M) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return m, nil }

func (m *M) View() string {
	style := m.Column().Style().Width(
		m.Column().Content - 2,
	).Border(lipgloss.NormalBorder()).Background(
		lipgloss.Color("5")) // DEBUG
	return m.RenderOrDie(
		zone.Mark(m.ID(), style.Render(m.name)),
	)
}

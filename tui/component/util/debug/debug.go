package debug

import (
	"github.com/charmbracelet/bubbletea"
	"github.com/minkezhang/truffle/tui/component/util/logger"

	l "github.com/minkezhang/truffle/tui/util/logger"
)

type M struct {
	logger tea.Model
	key    tea.Model
}

func Init() M {
	return M{
		logger: logger.Init(),
		key:    KeyLogger{},
	}
}

func (m M) Init() tea.Cmd { return nil }
func (m M) View() string  { return m.logger.View() }

func (m M) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.key.Update(msg)

	return m, nil
}

type KeyLogger struct{}

func (m KeyLogger) Init() tea.Cmd { return nil }
func (m KeyLogger) View() string  { return "" }

func (m KeyLogger) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		l.Debugf(
			"KeyInput { Type:%v, Runes:'%s' }",
			msg.Type,
			string(msg.Runes),
		)
	case tea.MouseMsg:
		l.Debugf(
			"MouseInput { X:%d, Y:%d, Action:%v, Button:%v }",
			msg.X,
			msg.Y,
			msg.Action,
			msg.Button,
		)
	}
	return m, nil
}

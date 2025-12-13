package debug

import (
	"github.com/76creates/stickers/flexbox"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/minkezhang/truffle/tui/component/util/logger"

	"github.com/minkezhang/truffle/tui/util/logging"
)

const (
	w = 100
)

type M struct {
	flex *flexbox.FlexBox
	c    *flexbox.Cell

	width int

	logger tea.Model
	key    tea.Model
}

func Init() M {
	c := flexbox.NewCell(1, 1).SetStyle(
		lipgloss.NewStyle().Background(lipgloss.Color("6")),
	).SetMinHeight(logging.Size()).SetMinWidth(w)
	f := flexbox.New(0, 0)
	f.AddRows(
		[]*flexbox.Row{
			f.NewRow().AddCells(c),
		},
	)
	return M{
		logger: logger.Init(logger.O{
			Width: w,
		}),
		key:    KeyLogger{},
		flex:   f,
		c:      c,
		width: w,
	}
}

func (m M) Init() tea.Cmd { return nil }
func (m M) View() string {
	m.c.SetContent(m.logger.View())
	return m.flex.Render()
}

func (m M) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.key.Update(msg)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if msg.Width >= m.width {
			m.flex.SetWidth(msg.Width)
		}
	}

	return m, nil
}

type KeyLogger struct{}

func (m KeyLogger) Init() tea.Cmd { return nil }
func (m KeyLogger) View() string  { return "" }

func (m KeyLogger) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		logging.Debugf(
			"KeyInput { Type:%v, Runes:'%s' }",
			msg.Type,
			string(msg.Runes),
		)
	case tea.MouseMsg:
		if msg.Action != tea.MouseActionMotion { // Reduce some spam
			logging.Debugf(
				"MouseInput { X:%d, Y:%d, Action:%v, Button:%v }",
				msg.X,
				msg.Y,
				msg.Action,
				msg.Button,
			)
		}
	case tea.WindowSizeMsg:
		logging.Debugf(
			"ResizeInput { W:%d, H:%d }",
			msg.Width,
			msg.Height,
		)
	}
	return m, nil
}

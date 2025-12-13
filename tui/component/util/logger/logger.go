package logger

import (
	"github.com/76creates/stickers/flexbox"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/minkezhang/truffle/tui/util/logging"
)

const (
	w = 80
)

var (
	background = map[logging.Severity]lipgloss.TerminalColor{
		logging.SeverityDebug:   lipgloss.Color("#000FFF"),
		logging.SeverityInfo:    lipgloss.Color("#0000FF"),
		logging.SeverityWarning: lipgloss.Color("#FFA500"),
		logging.SeverityError:   lipgloss.Color("#FF0000"),
	}
	foreground = map[logging.Severity]lipgloss.TerminalColor{
		logging.SeverityDebug:   lipgloss.Color("#FFFFFF"),
		logging.SeverityInfo:    lipgloss.Color("#FFFFFF"),
		logging.SeverityWarning: lipgloss.Color("#FFFF00"),
		logging.SeverityError:   lipgloss.Color("#FFFFFF"),
	}
)

type M struct {
	flex  *flexbox.FlexBox
	width int // min-width
}

type O struct {
	Width int // min-width
}

func Init(o O) M {
	flex := flexbox.New(0, 0)
	rows := []*flexbox.Row{}
	for _ = range logging.Size() {
		r := flex.NewRow()
		c := flexbox.NewCell(0, 0)
		c.SetContent("")
		r.AddCells(c)
		rows = append(rows, r)
	}
	flex.SetRows(rows)
	return M{
		flex:  flex,
		width: o.Width,
	}
}

func (m M) Init() tea.Cmd { return nil }

func (m M) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if msg.Width >= m.width {
			m.flex.SetWidth(msg.Width)
		}
	}
	return m, nil
}

func (m M) View() string {
	style := lipgloss.NewStyle()

	for i, msg := range logging.Messages() {
		r := m.flex.GetRow(logging.Size() - i - 1)
		r.SetStyle(style.Background(background[msg.S]))
		r.GetCell(0).SetContent(
			style.Foreground(foreground[msg.S]).Render(msg.String()),
		)
	}
	return m.flex.Render()
}

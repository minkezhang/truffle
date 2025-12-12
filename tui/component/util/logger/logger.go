package logger

import (
	"github.com/76creates/stickers/flexbox"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	l "github.com/minkezhang/truffle/tui/util/logger"
)

const (
	w = 100
)
var (
	colors = map[l.Severity]lipgloss.TerminalColor{
		l.SeverityDebug:   lipgloss.Color("#AAAAAA"),
		l.SeverityInfo:    lipgloss.Color("#FFFFFF"),
		l.SeverityWarning: lipgloss.Color("#FFFF00"),
		l.SeverityError:   lipgloss.Color("#FF0000"),
	}
)

type M struct {
	flex *flexbox.FlexBox
	width int
}

func Init() M {
	return M{
		flex: flexbox.New(0, 0).SetWidth(w),
		width: w,
	}
}

func (m M) Init() tea.Cmd { return nil }

func (m M) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	l.Debugf("%v", msg)
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if msg.Width >= w {
			m.flex.SetWidth(msg.Width)
		}
	}
	return m, nil
}

func (m M) View() string {
	style := lipgloss.NewStyle()

	rows := []*flexbox.Row{}
	for _, msg := range l.Messages() {
		r := m.flex.NewRow().AddCells(
			flexbox.NewCell(0, 0).SetContent(
				style.Foreground(colors[msg.S]).Background(lipgloss.Color("7")).Render(msg.String()),
			),
		)
		r.SetStyle(style.Background(lipgloss.Color("7")))
		rows = append(rows, r)
	}

	m.flex.SetRows(rows)
	return m.flex.Render()
}

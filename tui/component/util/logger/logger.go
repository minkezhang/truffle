package logger

import (
	"strings"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	l "github.com/minkezhang/truffle/tui/util/logger"
)

var (
	colors = map[l.Severity]lipgloss.TerminalColor{
		l.SeverityDebug:   lipgloss.Color("#AAAAAA"),
		l.SeverityInfo:    lipgloss.Color("#FFFFFF"),
		l.SeverityWarning: lipgloss.Color("#FFFF00"),
		l.SeverityError:   lipgloss.Color("#FF0000"),
	}
)

type M struct{}

func Init() M { return M{} }

func (m M) Init() tea.Cmd                           { return nil }
func (m M) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return m, nil }

func (m M) View() string {
	style := lipgloss.NewStyle().Background(lipgloss.Color("5"))
	var s strings.Builder

	messages := []string{}
	for _, m := range l.Messages() {
		messages = append(messages, style.Foreground(colors[m.S]).Render(m.String()))
	}

	for i := len(messages); i < l.Size(); i++ {
		messages = append([]string{""}, messages...)
	}
	s.WriteString(strings.Join(messages, "\n"))

	return s.String()
}

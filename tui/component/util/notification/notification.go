// TODO(minkezhang): Add timer.
// TODO(minkezhang): Add prefix.
package notification

import (
	"fmt"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"

	model_ui "github.com/minkezhang/truffle/tui/component/util/model"
)

const (
	ellipsis = "…"
)

type Type int

const (
	TypeNotice Type = iota
	TypeWarning
)

var (
	styles = map[Type]func(s lipgloss.Style) lipgloss.Style{
		TypeNotice: func(s lipgloss.Style) lipgloss.Style {
			return s.Background(lipgloss.Color("8"))
		},
		TypeWarning: func(s lipgloss.Style) lipgloss.Style {
			return s.Background(lipgloss.Color("3")).Foreground(lipgloss.Color("0")).Bold(true)
		},
	}

	prefix = map[Type]string{
		TypeNotice:  "NOTICE",
		TypeWarning: "WARNING",
	}
)

type O struct {
	model_ui.O
}

func New(o O) *M {
	return &M{
		Base: model_ui.New(o.O),
	}
}

type M struct {
	*model_ui.Base

	t       Type
	message string
}

func (m *M) Init() tea.Cmd { return nil }

func (m *M) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case model_ui.NoticeMsg:
		m.t = TypeNotice
		m.message = string(msg)
	case model_ui.WarningMsg:
		m.t = TypeWarning
		m.message = error(msg).Error()
	}
	return m, nil
}

func (m *M) body() string {
	message := ""
	if m.message != "" {
		message = fmt.Sprintf("%s: %s", prefix[m.t], m.message)
	}
	return ansi.Truncate(message, m.Column().Content, ellipsis)
}

func (m *M) View() string {
	style := styles[m.t](m.Column().Style()).MarginTop(1)
	return m.RenderOrDie(style.Render(m.body()))
}

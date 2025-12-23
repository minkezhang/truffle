package notification

import (
	"time"

	"github.com/charmbracelet/bubbles/timer"
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
			return s.Background(lipgloss.Color("3")).Foreground(lipgloss.Color("9")).Bold(true)
		},
	}
)

type O struct {
	model_ui.O
}

func New(o O) *M {
	return &M{
		Base:  model_ui.New(o.O),
		timer: timer.New(5 * time.Second),
	}
}

type message struct {
	t Type
	v string
}

type M struct {
	*model_ui.Base

	id int

	messages []message
	timer    timer.Model
}

const (
	timeout = 5 * time.Second
)

func (m *M) Init() tea.Cmd { return m.timer.Init() }

func (m *M) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case timer.TimeoutMsg:
		if m.timer.ID() == msg.ID {
			m.messages = m.messages[:len(m.messages)-1]
		}
		if len(m.messages) > 0 {
			m.timer.Timeout = timeout
			cmds = append(cmds, m.timer.Start())
		}
	case model_ui.NoticeMsg:
		m.timer.Timeout = timeout
		if !m.timer.Running() {
			cmds = append(cmds, m.timer.Start())
		}
		m.messages = append(m.messages, message{
			t: TypeNotice,
			v: string(msg),
		})
	case model_ui.WarningMsg:
		m.timer.Timeout = timeout
		if !m.timer.Running() {
			cmds = append(cmds, m.timer.Start())
		}
		m.messages = append(m.messages, message{
			t: TypeWarning,
			v: error(msg).Error(),
		})
	}

	var c tea.Cmd

	m.timer, c = m.timer.Update(msg)

	cmds = append(cmds, c)

	return m, tea.Batch(cmds...)
}

func (m *M) head() message {
	if len(m.messages) == 0 {
		return message{}
	}
	return m.messages[len(m.messages)-1]
}

func (m *M) body(n message) string {
	return ansi.Truncate(n.v, m.Column().Content, ellipsis)
}

func (m *M) View() string {
	if len(m.messages) == 0 {
		return ""
	}
	h := m.head()
	style := styles[h.t](m.Column().Style()).MarginBottom(1)
	return m.RenderOrDie(style.Render(m.timer.View() + m.body(h)))
}

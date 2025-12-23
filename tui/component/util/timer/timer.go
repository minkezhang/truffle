package timer

import (
	"time"

	"github.com/charmbracelet/bubbletea"
)

const (
	granularity = 500 * time.Millisecond
)

type Timeout struct {
	ID string
}

type tick struct {
	id       string
	enqueued time.Time
}

type start struct {
	id    string
	start time.Time
}

type pause tick

type stop struct {
	id string
}

type M struct {
	id string
	d  time.Duration

	running bool
	start   time.Time
	elapsed time.Duration

	now   func() time.Time
	sleep func(d time.Duration)
}

func New(d time.Duration) *M {
	return &M{
		id:    "foo",
		d:     d,
		now:   time.Now,
		sleep: time.Sleep,
	}
}

func (m *M) Init() tea.Cmd { return nil }

func (m *M) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case start:
		if !m.running && m.id == msg.id {
			m.running = true
			m.start = msg.start
			enqueued := m.now()
			cmds = append(cmds, func() tea.Msg {
				return tick{
					id:       m.id,
					enqueued: enqueued,
				}
			})
		}
	case stop:
		if m.id == msg.id {
			m.running = false
			m.elapsed = 0
		}
	case pause:
		if m.id == msg.id {
			m.running = false
			m.elapsed += m.now().Sub(msg.enqueued)
			if m.elapsed >= m.d && m.running {
				cmds = append(cmds, func() tea.Msg {
					return Timeout{
						ID: m.id,
					}
				})
			}
		}
	case tick:
		if m.running && m.id == msg.id {
			m.elapsed += m.now().Sub(msg.enqueued)
			if m.elapsed >= m.d {
				if m.running {
					m.running = false
					cmds = append(cmds, func() tea.Msg {
						return Timeout{
							ID: m.id,
						}
					})
				}
			} else {
				enqueued := m.now()
				cmds = append(cmds, func() tea.Msg {
					m.sleep(time.Second)
					return tick{
						id:       m.id,
						enqueued: enqueued,
					}
				})
			}
		}
	}

	return m, tea.Batch(cmds...)
}

func (m *M) View() string { return m.elapsed.String() }

func (m *M) Pause() tea.Cmd {
	enqueued := m.now()
	return func() tea.Msg {
		return pause{
			id:       m.id,
			enqueued: enqueued,
		}
	}
}

func (m *M) Stop() tea.Cmd {
	return func() tea.Msg {
		return stop{
			id: m.id,
		}
	}
}

func (m *M) Start() tea.Cmd {
	enqueued := m.now()
	return func() tea.Msg {
		return start{
			id:    m.id,
			start: enqueued,
		}
	}
}

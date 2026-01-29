package timer

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbletea"
	"github.com/minkezhang/truffle-api/util/generator"
	"github.com/minkezhang/truffle/tui/component/directory/base"
	"github.com/minkezhang/truffle/tui/component/focusable"
)

type TimerMessageType int

const (
	TimerMessageStart TimerMessageType = iota
	TimerMessageStop
	TimerMessageInterrupt
)

var (
	g = generator.New(generator.O{N: 16})
)

type O struct {
	ParentID string
}

type Node struct {
	*focusable.Node

	run_id      string
	last_update time.Time
	start       time.Time
	duration    time.Duration
}

func New(o O) *Node {
	return &Node{
		Node: focusable.New("", o.ParentID, 0),
	}
}

func (n *Node) Init() tea.Cmd {
	return func() tea.Msg {
		return base.RegisterMessage{
			Node: n,
		}
	}
}

type tick_message struct {
	id     string // Node ID
	run_id string
	sleep  time.Duration
}

type TimerMessage struct {
	ID    string
	RunID string
	Type  TimerMessageType
}

func (n *Node) do_tick(run_id string, sleep time.Duration) tea.Cmd {
	return func() tea.Msg {
		return tick_message{
			id:     n.ID(),
			run_id: run_id,
			sleep:  sleep,
		}
	}
}

func (n *Node) Start(duration time.Duration) tea.Cmd {
	var cmds []tea.Cmd

	run_id := g.Generate()

	if d := time.Now().Sub(n.start); d < n.duration {
		cmds = append(cmds, func() tea.Msg {
			return TimerMessage{
				ID:    n.ID(),
				RunID: n.run_id,
				Type:  TimerMessageInterrupt,
			}
		})
	}

	n.run_id = run_id
	n.start = time.Now()
	n.duration = duration

	cmds = append(cmds,
		n.do_tick(run_id, 100*time.Millisecond),
		func() tea.Msg {
			return TimerMessage{
				ID:    n.ID(),
				RunID: run_id,
				Type:  TimerMessageStart,
			}
		},
	)

	return tea.Batch(cmds...)
}

func (n *Node) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tick_message:
		if msg.id == n.ID() {
			// Mismatching run_id implies a previous iteration was
			// still running. We have handled sending
			// TimerMessageInterrupt during n.Start().
			if n.run_id == msg.run_id {
				if d := time.Now().Sub(n.start); d >= n.duration {
					cmds = append(cmds, func() tea.Msg {
						return TimerMessage{
							ID:    n.ID(),
							Type:  TimerMessageStop,
							RunID: msg.run_id,
						}
					})
				} else {
					cmds = append(cmds, n.do_tick(msg.run_id, msg.sleep))
				}
			}
		}
	}
	return n, tea.Batch(cmds...)
}

func (n *Node) View() string {
	return fmt.Sprintf(
		"ID: %v, run_id: %v, time remaining: %v",
		n.ID(),
		n.run_id,
		n.start.Add(n.duration).Sub(time.Now()),
	)
}

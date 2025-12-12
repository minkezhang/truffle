package logger

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

type Severity int

func (s Severity) String() string {
	return map[Severity]string{
		SeverityDebug:   "D",
		SeverityInfo:    "I",
		SeverityWarning: "W",
		SeverityError:   "E",
	}[s]
}

const (
	SeverityDebug Severity = iota
	SeverityInfo
	SeverityWarning
	SeverityError
)

var (
	buf = New(O{
		Size: 5,
	})
)

type t func() time.Time

type O struct {
	Size int
	now  t
}

type L struct {
	lines []M
	size  int
	tail  int
	count int
	now   t
	mux   sync.Mutex
}

type M struct {
	S Severity
	M string
	T time.Time
}

func (m M) String() string {
	return fmt.Sprintf(
		"%v%v: %v",
		m.S.String(),
		m.T.Format("2006-01-02T15:04:05 -070000"),
		m.M,
	)
}

func New(o O) *L {
	l := &L{
		lines: make([]M, o.Size),
		size:  o.Size,
		tail:  0,
		count: 0,
		now:   time.Now,
	}
	if o.now != nil {
		l.now = o.now
	}
	return l
}

func (l *L) push(s Severity, m string) {
	l.mux.Lock()
	defer l.mux.Unlock()

	ms := strings.Split(strings.ReplaceAll(m, "\r\n", "\n"), "\n")
	t := l.now()
	for _, m := range ms {
		l.lines[l.tail] = M{
			S: s,
			M: m,
			T: t,
		}
		l.tail = (l.tail + 1) % l.size

		if l.count < l.size {
			l.count += 1
		}
	}
}

func (l *L) Debug(v string)   { l.push(SeverityDebug, v) }
func (l *L) Info(v string)    { l.push(SeverityInfo, v) }
func (l *L) Warning(v string) { l.push(SeverityWarning, v) }
func (l *L) Error(v string)   { l.push(SeverityError, v) }

func (l *L) Debugf(f string, args ...any)   { l.push(SeverityDebug, fmt.Sprintf(f, args...)) }
func (l *L) Infof(f string, args ...any)    { l.push(SeverityInfo, fmt.Sprintf(f, args...)) }
func (l *L) Warningf(f string, args ...any) { l.push(SeverityWarning, fmt.Sprintf(f, args...)) }
func (l *L) Errorf(f string, args ...any)   { l.push(SeverityError, fmt.Sprintf(f, args...)) }

func (l *L) Messages() []M {
	l.mux.Lock()
	defer l.mux.Unlock()

	res := make([]M, 0, l.size)

	for i := 0; i < l.count; i++ {
		index := (l.tail + l.size - l.count + i) % l.size
		res = append(res, l.lines[index])
	}

	return res
}

func (l *L) Size() int { return l.size }

func (l *L) Len() int {
	l.mux.Lock()
	defer l.mux.Unlock()

	return l.count
}

func Debug(v string)   { buf.push(SeverityDebug, v) }
func Info(v string)    { buf.push(SeverityInfo, v) }
func Warning(v string) { buf.push(SeverityWarning, v) }
func Error(v string)   { buf.push(SeverityError, v) }

func Debugf(f string, args ...any)   { buf.push(SeverityDebug, fmt.Sprintf(f, args...)) }
func Infof(f string, args ...any)    { buf.push(SeverityInfo, fmt.Sprintf(f, args...)) }
func Warningf(f string, args ...any) { buf.push(SeverityWarning, fmt.Sprintf(f, args...)) }
func Errorf(f string, args ...any)   { buf.push(SeverityError, fmt.Sprintf(f, args...)) }

func Messages() []M { return buf.Messages() }
func Size() int     { return buf.size }

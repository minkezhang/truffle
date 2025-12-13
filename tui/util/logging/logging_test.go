package logging

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestPush(t *testing.T) {
	t.Run("NoOverwrite", func(t *testing.T) {
		l := New(O{
			Size: 2,
			now:  func() time.Time { return time.Time{} },
		})

		l.Debug("0")

		want := []M{
			M{
				S: SeverityDebug,
				M: "0",
			},
		}

		got := l.Messages()
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("Messages() mismatch (-want +got):\n%v", diff)
		}
	})
	t.Run("Overwrite", func(t *testing.T) {
		l := New(O{Size: 2})
		l.now = func() time.Time { return time.Time{} }

		l.Debug("0")
		l.Debug("1")
		l.Debug("2")
		l.Debug("3")

		want := []M{
			M{
				S: SeverityDebug,
				M: "2",
			},
			M{
				S: SeverityDebug,
				M: "3",
			},
		}

		got := l.Messages()
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("Messages() mismatch (-want +got):\n%v", diff)
		}
	})
	t.Run("Overwrite/MultiLine", func(t *testing.T) {
		l := New(O{Size: 2})
		l.now = func() time.Time { return time.Time{} }

		l.Debug("0\n1\n2\n3")

		want := []M{
			M{
				S: SeverityDebug,
				M: "2",
			},
			M{
				S: SeverityDebug,
				M: "3",
			},
		}

		got := l.Messages()
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("Messages() mismatch (-want +got):\n%v", diff)
		}
	})
}

package render

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func Validate(s string, max_width int) error {
	w := lipgloss.Width(s)
	if w <= max_width {
		return nil
	}
	return fmt.Errorf(
		"rendered string exceeded bounding box: %d > %d\n  %s\n%s\n",
		w,
		max_width,
		strings.Repeat(" ", max_width)+"↓",
		with_xy_axis(s),
	)
}

func with_xy_axis(s string) string {
	w := lipgloss.Width(s)
	var axis string
	if w < 10 {
		axis = []string{
			"",
			"│",
			"├─",
			"├──",
			"├───",
			"├────",
			"├────┐",
			"├────┬─",
			"├────┬──",
			"├────┬───",
		}[w]
	} else {
		axis = fmt.Sprintf(
			"%v%v%v",
			"├────┬────",
			strings.Repeat(
				"┼────┬────",
				(w/10)-1,
			),
			[]string{
				"",
				"┤",
				"┼─",
				"┼──",
				"┼───",
				"┼────",
				"┼────┐",
				"┼────┬─",
				"┼────┬──",
				"┼────┬───",
			}[w%10],
		)
	}
	return fmt.Sprintf("  %v\n%v", axis, with_y_axis(s))

}

func with_y_axis(s string) string {
	h := lipgloss.Height(s)
	var axis string
	if h < 10 {
		axis = []string{
			"",
			"─",
			"┬│",
			"┬││",
			"┬│││",
			"┬││││",
			"┬││││└",
			"┬││││├│",
			"┬││││├││",
			"┬││││├│││",
		}[h]
	} else {
		axis = fmt.Sprintf(
			"%v%v%v",
			"┬││││├││││",
			strings.Repeat(
				"┼││││├││││",
				(h/10)-1,
			),
			[]string{
				"",
				"┴",
				"┼│",
				"┼││",
				"┼│││",
				"┼││││",
				"┼││││└",
				"┼││││├│",
				"┼││││├││",
				"┼││││├│││",
			}[h%10],
		)
	}
	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		strings.Join(strings.Split(axis, ""), "\n"),
		" ",
		s,
	)
}

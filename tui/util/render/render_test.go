package render

import (
	"testing"
)

func TestCheck(t *testing.T) {
	configs := []struct {
		name      string
		s         string
		max_width int
		success   bool
	}{
		{
			name:      "Trivial",
			s:         "",
			max_width: 0,
			success:   true,
		},
		{
			name:      "Trivial/Fail",
			s:         " ",
			max_width: 0,
			success:   false,
		},
		{
			name:      "Simple/ZeroWidth",
			s:         "The\nfox",
			max_width: 3,
			success:   true,
		},
		{
			name:      "Simple/Fail/SecondLine",
			s:         "The quick brown fox\njumped over the lazy dog",
			max_width: len("The quick brown fox"),
			success:   false,
		},
		{
			name: "Simple/Pass",
			s: "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor\n" +
				"incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis\n" +
				"nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.\n" +
				"Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu\n" +
				"fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in\n" +
				"culpa qui officia deserunt mollit anim id est laborum.",
			max_width: 80,
			success:   true,
		},
	}

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			err := Validate(c.s, c.max_width)
			if err != nil && c.success {
				t.Errorf("check() returned a non-nil error: %v", err)
			} else if err == nil && !c.success {
				t.Errorf("check() unexpectedly succeeded")
			}
		})
	}
}

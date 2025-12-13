package main

// These imports will be used later on the tutorial. If you save the file
// now, Go might complain they are unused, but that's fine.
// You may also need to run `go mod tidy` to download bubbletea and its
// dependencies.
import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbletea"
	"github.com/lrstanley/bubblezone"
	// "github.com/minkezhang/truffle/tui/component/root"
	"github.com/minkezhang/truffle/tui/util/logging"

	"github.com/minkezhang/truffle/tui/component/util/logger"
	// "github.com/minkezhang/truffle/tui/component/util/debug"
)

func main() {
	logging.SetSize(10)

	// See https://github.com/lrstanley/bubblezone for more information.
	zone.NewGlobal()

	p := tea.NewProgram(
		logger.Init(logger.O{Width: 80}),
		// debug.Init(),
		// root.New(),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Run() returned unexpected error: %v", err)
		os.Exit(1)
	}
}

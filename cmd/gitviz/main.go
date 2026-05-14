package main

import (
	"fmt"
	"os"

	"github.com/alikhaleghi/gitviz/internal/tui"
)

func main() {
	if err := tui.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "gitviz failed: %v\n", err)
		os.Exit(1)
	}
}

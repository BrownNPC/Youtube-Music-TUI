//util functions for playlists table

package main

import (
	"github.com/charmbracelet/bubbles/table"
)

func BuildTPlaylists(WindowWidth int, WwindowHeight int) table.Model {
	t := table.New(
		table.WithFocused(true),
		table.WithWidth(PercentageOf(WindowWidth, 20)),
		table.WithHeight(WwindowHeight-2),
	)

	return t
}

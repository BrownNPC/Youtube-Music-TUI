//util functions for playlists table

package main

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

func BuildTPlaylists(WindowWidth int, WindowHeight int, Playlists []playlist) table.Model {

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("white")).
		Background(lipgloss.Color("57")).
		Bold(false)

	t := table.New(
		table.WithStyles(s),
		table.WithFocused(true),
	)

	t.SetWidth(PercentageOf(WindowWidth, 25))
	t.SetHeight(WindowHeight - 6)
	t.SetColumns([]table.Column{
		{Title: "Playlists", Width: t.Width() - 2},
	})

	t.SetRows(TPlaylistsBuildRows(Playlists))

	return t
}

func TPlaylistsBuildRows(Playlists []playlist) []table.Row {
	var rows []table.Row
	for _, p := range Playlists {
		rows = append(rows, table.Row{p.Title})
	}
	return rows
}

//util functions for tracks table

package main

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

func BuildTTracks(WindowWidth int, WindowHeight int, p playlist) table.Model {

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
		table.WithFocused(false),
	)

	t.SetWidth(PercentageOf(WindowWidth, 72))
	t.SetHeight(WindowHeight - 6)
	t.SetColumns([]table.Column{
		{Title: "Name", Width: PercentageOf(t.Width(), 80) - 2},
		{Title: "Channel", Width: PercentageOf(t.Width(), 20) - 2},
	})

	t.SetRows(TTracksBuildRows(p))

	return t
}

func TTracksBuildRows(p playlist) []table.Row {
	var rows []table.Row
	for _, p := range p.Entries {
		rows = append(rows, table.Row{p.Title, p.Channel})

	}
	return rows
}

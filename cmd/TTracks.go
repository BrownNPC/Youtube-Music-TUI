//util functions for playlists table

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
		table.WithFocused(true),
	)

	t.SetWidth(PercentageOf(WindowWidth, 72))
	t.SetHeight(WindowHeight - 7)
	t.SetColumns([]table.Column{
		{Title: "Name", Width: PercentageOf(t.Width(), 80) - 2},
		{Title: "Artist", Width: PercentageOf(t.Width(), 20) - 2},
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

func (m *model) refreshTracks() {
	m.TTracks = BuildTTracks(m.sizeX, m.sizeY, m.Playlists[m.Cursor])
}

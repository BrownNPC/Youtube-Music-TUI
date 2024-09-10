package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	Playlists []playlist

	TPlaylists table.Model
	TTracks    table.Model

	sizeX, sizeY int

	output string

	Cursor int
}

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("240"))

func (m model) View() string {
	// Return a string representation of the model's view

	var tabs = lipgloss.JoinHorizontal(
		0, baseStyle.Render(m.TPlaylists.View()),
		baseStyle.Render(m.TTracks.View()),
	)

	return lipgloss.JoinVertical(0, tabs, m.output)

	// return baseStyle.Render(m.TPlaylists.View())
}

func (m model) Init() tea.Cmd {
	// Initialize any commands or leave it nil if not needed
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?

		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		case "enter":
			m.Cursor = m.TPlaylists.Cursor()
			m.refreshTracks()

		}
		var cmd tea.Cmd
		m.TPlaylists, cmd = m.TPlaylists.Update(msg)

		m.TTracks, _ = m.TTracks.Update(msg)

		m.output = fmt.Sprint(m.Playlists[m.Cursor].Entries[m.TTracks.Cursor()].Channel)
		return m, cmd
	case tea.WindowSizeMsg:
		m.sizeX, m.sizeY = msg.Width, msg.Height
		m.refreshPlaylists()
		m.refreshTracks()

	}

	return m, nil
}

func main() {
	m := model{}
	fmt.Println("Updating playlist cache, this is a one time operation...")
	m.Playlists = append(m.Playlists, QuickLoadPlaylist("PLkcA3mJSVisBLbLhQ6ZnTCi9nGHTVUDaI"))
	m.Playlists = append(m.Playlists, QuickLoadPlaylist("PLkcA3mJSVisCozQtw7xVXn_zPzrjWsvr9"))

	program := tea.NewProgram(m, tea.WithAltScreen())
	program.Run()
	// for _, e := range m.Playlists[0].Entries {
	// 	fmt.Println(e.Title)
	// }

}

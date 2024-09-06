package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	Playlists []playlist

	TPlaylists table.Model
}

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

func (m model) View() string {
	// Return a string representation of the model's view

	return baseStyle.Render(m.TPlaylists.View()) // You can add more detailed output here based on the playlist
}

func (m model) Init() tea.Cmd {
	// Initialize any commands or leave it nil if not needed
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle incoming messages and update the model accordingly
	// For now, just return the model unchanged
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "down":
			m.TPlaylists.MoveDown(1)

		case "up":
			m.TPlaylists.MoveUp(1)
		case "ctrl+c":
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.TPlaylists = BuildTPlaylists(msg.Width, msg.Height)
	}

	return m, nil
}

func main() {
	m := model{}

	p := QuickLoadPlaylist("PLkcA3mJSVisBLbLhQ6ZnTCi9nGHTVUDaI")
	m.Playlists = append(m.Playlists, p)

	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		fmt.Printf("could not start program: %s\n", err)
		os.Exit(1)
	}
}
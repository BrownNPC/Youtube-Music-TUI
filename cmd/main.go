package main

import (
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

func boxStyle(width int, height int, bg lipgloss.Color) lipgloss.Style {
	return lipgloss.NewStyle().
		Background(bg).
		Foreground(lipgloss.Color("0")).
		Width(width).
		Height(height).
		Align(lipgloss.Center, lipgloss.Center)
}

func (m model) View() string {
	// Return a string representation of the model's view

	return baseStyle.Render(m.TPlaylists.View())
}

func (m model) Init() tea.Cmd {
	// Initialize any commands or leave it nil if not needed
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle incoming messages and update the model accordingly
	// For now, just return the model unchanged
	switch msg := msg.(type) {

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

	program := tea.NewProgram(m, tea.WithAltScreen())
	program.Run()

}

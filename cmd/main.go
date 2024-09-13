package main

import (
	_ "embed"
	"sync"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	mpv "github.com/gen2brain/go-mpv"
)

type model struct {
	Playlists []playlist

	TPlaylists table.Model
	TTracks    table.Model

	sizeX, sizeY int // window size

	Cursor   int
	Shuffled bool

	Player *mpv.Mpv
}

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("240"))

var activeStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("57"))

func (m model) View() string {
	// Return a string representation of the model's view
	var tabs string

	if m.TPlaylists.Focused() {

		tabs = lipgloss.JoinHorizontal(
			0, activeStyle.Render(m.TPlaylists.View()),
			baseStyle.Render(m.TTracks.View()),
		)
	} else if m.TTracks.Focused() {
		tabs = lipgloss.JoinHorizontal(
			0, baseStyle.Render(m.TPlaylists.View()),
			activeStyle.Render(m.TTracks.View()),
		)
	}
	return lipgloss.JoinVertical(0, tabs, "bruh")

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
			if m.TPlaylists.Focused() {
				m.Cursor = m.TPlaylists.Cursor()
				m.refreshTracks()
				m.swapFocus()
			}
		case "tab":
			m.swapFocus()
		}
		var cmd tea.Cmd
		m.TPlaylists, cmd = m.TPlaylists.Update(msg)
		m.TTracks, _ = m.TTracks.Update(msg)
		return m, cmd
	case tea.WindowSizeMsg:
		m.sizeX, m.sizeY = msg.Width, msg.Height
		m.refreshPlaylists()
		m.refreshTracks()
	}
	return m, nil
}

func (m *model) swapFocus() {
	if m.TPlaylists.Focused() {
		m.TPlaylists.Blur()
		m.TTracks.Focus()
	} else if m.TTracks.Focused() {
		m.TTracks.Blur()
		m.TPlaylists.Focus()
	}
}
func main() {
	m := model{}

	handleCommandLineArgs() // api.go
	cfg, err := LoadConfig()
	if err != nil {
		panic(err)
	}

	wg := sync.WaitGroup{}
	for _, id := range cfg.IDs {
		wg.Add(1)
		go func(id string) {
			defer wg.Done()
			m.Playlists = append(m.Playlists, QuickLoadPlaylist(id))
		}(id)
	}
	wg.Wait()
	program := tea.NewProgram(m, tea.WithAltScreen())
	program.Run()

}

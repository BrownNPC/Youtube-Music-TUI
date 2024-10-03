package main

import (
	_ "embed"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gen2brain/go-mpv"

	"github.com/charmbracelet/bubbles/progress"
)

type model struct {
	Playlists []playlist

	TPlaylists table.Model
	TTracks    table.Model

	sizeX, sizeY int // window size

	Cursor            int
	Shuffled          bool
	CurrentTrackIndex int

	Player      *mpv.Mpv
	Duration    string // format duration OSD (On-Screen Display)
	ProgressBar progress.Model

	Playing       bool
	TrackPlaying  string // currently playing track
	VolumePercent int

	NextTrackStream  string // stream url for the next track
	ShuffledOrder    []int  // shuffled order of tracks (indexes)
	LoadingNextTrack bool   // used internally
}

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("240"))

var activeStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("57"))

type TickMsg time.Time

func Tick() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

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

	return lipgloss.JoinVertical(0, tabs,
		m.TrackPlaying, m.Duration,
	)

}

func (m model) Init() tea.Cmd {
	return Tick() // Start the ticker
}

func (m *model) HandleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit

	case "right", "left":

	case "enter":
		if m.TPlaylists.Focused() {
			m.Cursor = m.TPlaylists.Cursor()
			m.refreshTracks()
			m.swapFocus()
			return m, nil
		}
		if m.TTracks.Focused() {
			m.CurrentTrackIndex = m.TTracks.Cursor()
		}
	case "tab":
		m.swapFocus()
		return m, nil

	case " ":
		//play/pause
	case ",", ".":
		//volume up and down

		return m, nil

	case "s":
		//shuffle
		fmt.Println(m.ShuffledOrder)
	}
	var cmd tea.Cmd
	m.TPlaylists, cmd = m.TPlaylists.Update(msg)
	m.TTracks, _ = m.TTracks.Update(msg)

	return m, cmd
}

func (m *model) getNextTrack() {
	var TrackIndex = m.CurrentTrackIndex

	if m.Shuffled {
		TrackIndex = m.ShuffledOrder[m.ShuffledOrder[m.CurrentTrackIndex]+1%len(m.ShuffledOrder)]
	} else {
		TrackIndex = m.CurrentTrackIndex + 1%len(m.TTracks.Rows())
	}
	m.CurrentTrackIndex = TrackIndex
	url := m.Playlists[m.Cursor].Entries[TrackIndex].Url
	streamurl, err := getYoutubeStreamURL(url)
	if err != nil {
		log.Fatal(err)
	}
	m.TTracks.SetCursor(TrackIndex)
	fmt.Print(m.Playlists[m.Cursor].Entries[TrackIndex].Title)
	if err != nil {
		log.Fatal(err)
	}
	if streamurl == "" {
		panic("streamurl is empty")
	}
	m.NextTrackStream = streamurl

}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// Is it a key press?
	case tea.KeyMsg:
		// Cool, what was the actual key pressed?

		return m.HandleKeyPress(msg)

	case TickMsg:

		m.UpdateDuration()

		return m, Tick()
	case tea.WindowSizeMsg:
		m.sizeX, m.sizeY = msg.Width, msg.Height
		m.refreshPlaylists()
		m.refreshTracks()
		m.ProgressBar.Width = m.sizeX - 10

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

	handleCommandLineArgs() // api.go
	cfg, err := LoadConfig()
	if err != nil {
		panic(err)
	}
	plr := mpv.New()
	if err != nil {
		panic(err)
	}
	err = plr.Initialize()
	if err != nil {
		panic(err)
	}
	m := model{Player: plr,
		ProgressBar: progress.New(progress.WithScaledGradient("#FF7CCB", "#FDFF8C"),
			progress.WithoutPercentage()),
		TrackPlaying: "Nothing is playing"}

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

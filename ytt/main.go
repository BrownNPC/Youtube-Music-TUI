package main

import (
	_ "embed"
	"fmt"
	"log"
	"math/rand/v2"
	"sync"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gen2brain/go-mpv"

	"github.com/charmbracelet/bubbles/progress"
)

var nextTrackFetched bool = false
var CurrentTrackIndex int = 0
var RealTrackIndex int = 0 // used for getting track name
var CancelNextTrackFetch bool = false

type model struct {
	Playlists []playlist

	TPlaylists table.Model
	TTracks    table.Model

	sizeX, sizeY int // window size

	Cursor   int
	Shuffled bool

	Player      *mpv.Mpv
	Duration    string // format duration OSD (On-Screen Display)
	ProgressBar progress.Model

	Paused        bool
	TrackPlaying  string // currently playing track
	VolumePercent int

	NextTrackStream string  // stream url for the next track
	ProgressPercent float64 // 0.0 - 1.0
	ShuffledOrder   []int   // shuffled order of tracks (indexes)
	ChanNextTrack   chan string
}

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("240"))

var activeStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("57"))

func main() {

	handleCommandLineArgs() // api.go
	cfg, err := LoadConfig()
	if err != nil {
		panic(err)
	}
	plr := mpv.New()
	err = plr.Initialize()
	if err != nil {
		panic(err)
	}
	plr.RequestLogMessages("info")
	plr.SetProperty("pause", mpv.FormatFlag, true)
	m := model{Player: plr,
		ProgressBar: progress.New(progress.WithScaledGradient("#FF7CCB", "#FDFF8C"),
			progress.WithoutPercentage()),
		TrackPlaying: "Nothing is playing", ChanNextTrack: make(chan string)}

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

type TickMsg time.Time

func Tick() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

func (m model) Init() tea.Cmd {
	return Tick() // Start the ticker
}

func (m model) View() string {
	var tabs string
	// outline around active tab
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
	playpauseicon := "⏸"
	if m.Paused {
		playpauseicon = "▶"
	}

	shuffleicon := " "
	if m.Shuffled {
		shuffleicon = "🔀"
	}

	return lipgloss.JoinVertical(0, tabs,
		m.TrackPlaying, m.Duration,
		lipgloss.JoinHorizontal(lipgloss.Center,
			playpauseicon, "  ", m.ProgressBar.ViewAs(m.ProgressPercent), fmt.Sprint(m.VolumePercent, "%"), " ", shuffleicon),
	)

}

func (m *model) HandleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit

	case "right", "left":
		if msg.String() == "right" {
			m.Player.Command([]string{"seek", "10", "relative"})
		} else if msg.String() == "left" {
			m.Player.Command([]string{"seek", "-10", "relative"})
		}

	case "enter":
		// if table of playlists is focus, load the tracks from the selected playlist
		// and swap focus to the table of tracks
		if m.TPlaylists.Focused() {
			m.Cursor = m.TPlaylists.Cursor()
			m.refreshTracks()
			m.swapFocus()
			return m, nil
		}
		//if table of tracks is focused, play the selected track
		if m.TTracks.Focused() {
			CurrentTrackIndex = m.TTracks.Cursor()
			RealTrackIndex = CurrentTrackIndex
			m.generateShuffleOrder()
			go func() {
				// stop next track from being fetched since the user has selected a track on his own
				CancelNextTrackFetch = true
				url, err := getYoutubeStreamURL(m.Playlists[m.Cursor].Entries[CurrentTrackIndex].Url)
				if err != nil {
					fmt.Println("error getting stream url", err)
				}
				err = m.Player.Command([]string{"loadfile", url})
				if err != nil {
					log.Panicln("error loading track \n ", err, url)
				}
				m.Player.SetProperty("pause", mpv.FormatFlag, false)

			}()

		}
	case "tab":
		m.swapFocus()
		return m, nil
	case " ": //space bar
		m.Player.SetProperty("pause", mpv.FormatFlag, !m.Paused)
		return m, nil
	case ",", ".":
		//volume up and down
		if msg.String() == "," && m.VolumePercent >= 0 {
			m.Player.SetProperty("volume", mpv.FormatInt64, m.VolumePercent-10)
		} else if msg.String() == "." && m.VolumePercent <= 90 {
			m.Player.SetProperty("volume", mpv.FormatInt64, m.VolumePercent+10)
		}

		return m, nil

	case "s":
		//shuffle
		m.Shuffled = !m.Shuffled

	}
	var cmd tea.Cmd
	m.TPlaylists, cmd = m.TPlaylists.Update(msg)
	m.TTracks, _ = m.TTracks.Update(msg)

	return m, cmd
}
func (m *model) getNextTrack() string {
	var TrackIndex int
	CancelNextTrackFetch = false // reset
	if m.Shuffled {
		TrackIndex = (m.ShuffledOrder[CurrentTrackIndex] + 1) % len(m.ShuffledOrder)
	} else {
		TrackIndex = (CurrentTrackIndex + 1) % len(m.TTracks.Rows())
	}
	url := m.Playlists[m.Cursor].Entries[TrackIndex].Url
	streamurl, err := getYoutubeStreamURL(url)
	if err != nil {
		log.Fatal(err)
	}

	if streamurl == "" {
		panic("streamurl is empty")
	}
	// if the user played a track on his own while we were fetching the next track
	if CancelNextTrackFetch {
		CancelNextTrackFetch = false
		return ""
	}
	RealTrackIndex = TrackIndex
	return streamurl
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// Is it a key press?
	case tea.KeyMsg:
		return m.HandleKeyPress(msg)

	case TickMsg:
		// update volume percentage
		var newVolume interface{}
		newVolume, err := m.Player.GetProperty("volume", mpv.FormatInt64)
		if err != nil {
			log.Fatal(err)
		}
		m.VolumePercent = int(newVolume.(int64))

		// update paused
		var paused interface{}
		paused, err = m.Player.GetProperty("pause", mpv.FormatFlag)
		if err != nil {
			log.Fatal(err)
		}
		m.Paused = paused.(bool)

		// update current playing track
		if !m.Paused {
			m.TrackPlaying = m.Playlists[m.Cursor].Entries[RealTrackIndex].Title
		}

		//update track duration and progressbar percentage
		if !m.Paused {
			currentPlaybackTime, _ := m.Player.GetProperty("time-pos", mpv.FormatInt64)
			duration, err := m.Player.GetProperty("percent-pos", mpv.FormatInt64)
			if err != nil {
				duration = int64(0)
				currentPlaybackTime = int64(0)
			}
			m.ProgressPercent = float64(duration.(int64)) / float64(100)
			m.Duration = formatTime(currentPlaybackTime.(int64)) + " - " + formatTime(m.CurrentTrack().Duration)
		}

		// load next track if we are 90% of the way to the end

		if m.ProgressPercent >= 0.9 && !m.Paused && !nextTrackFetched {
			go func() {
				nextTrackFetched = true
				m.ChanNextTrack <- m.getNextTrack()
			}()
		}
		go func() {
			e := m.Player.WaitEvent(1)
			switch e.EventID {
			case mpv.EventEnd:
				if e.EndFile().Reason.String() == "eof" {
					if nextTrackFetched {
						m.NextTrackStream = <-m.ChanNextTrack // next track fetched
						if m.NextTrackStream == "" {
							return
						}
						err = m.Player.Command([]string{"loadfile", m.NextTrackStream})
						if err != nil {
							log.Fatal(err)
						}
						CurrentTrackIndex = CurrentTrackIndex + 1%len(m.TTracks.Rows())
						nextTrackFetched = false
					}
				}
			}
		}()
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

func (m *model) CurrentTrack() *Entry {
	return &m.Playlists[m.Cursor].Entries[CurrentTrackIndex]
}

func (m *model) UpdateDuration() {

	duration, err := m.Player.GetProperty("duration", mpv.FormatInt64)
	if err != nil {
		log.Fatal(err)
	}
	m.Duration = formatTime(int(duration.(int64)))

	track, err := m.Player.GetProperty("time-pos", mpv.FormatInt64)
	if err != nil {
		log.Fatal(err)
	}
	m.ProgressBar.SetPercent(float64(track.(int64)) / float64(duration.(int64)))
}

func (m *model) refreshTracks() {
	if len(m.Playlists) == 0 {
		m.TTracks = BuildTTracks(m.sizeX, m.sizeY, playlist{})
		return
	}
	m.TTracks = BuildTTracks(m.sizeX, m.sizeY, m.Playlists[m.Cursor])

}

func (m *model) generateShuffleOrder() {
	// create shuffle order
	m.ShuffledOrder = make([]int, len(m.TTracks.Rows()))
	for i := 0; i < len(m.TTracks.Rows()); i++ {
		m.ShuffledOrder[i] = i
	}

	// shuffle the order
	rand.Shuffle(len(m.ShuffledOrder), func(i, j int) {
		m.ShuffledOrder[i], m.ShuffledOrder[j] = m.ShuffledOrder[j], m.ShuffledOrder[i]
	})
}

func (m *model) refreshPlaylists() {
	m.TPlaylists = BuildTPlaylists(m.sizeX, m.sizeY, m.Playlists)
}

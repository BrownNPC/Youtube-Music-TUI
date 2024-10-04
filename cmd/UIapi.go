package main

import (
	"log"
	"math/rand/v2"

	"github.com/gen2brain/go-mpv"
)

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

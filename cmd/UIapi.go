package main

import (
	"math/rand/v2"
)

func (m *model) UpdateDuration() {
}

func (m *model) refreshTracks() {
	if len(m.Playlists) == 0 {
		m.TTracks = BuildTTracks(m.sizeX, m.sizeY, playlist{})
		return
	}
	m.TTracks = BuildTTracks(m.sizeX, m.sizeY, m.Playlists[m.Cursor])
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

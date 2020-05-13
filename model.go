package main

type Album struct {
	Id     int     `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Tracks []Track `json:"tracks"`
}

type Track struct {
	Id          int         `json:"id"`
	Title       string      `json:"title"`
	Artist      string      `json:"artist"`
	TrackNumber TrackNumber `json:"trackNumber"`
	Album       string      `json:"album"`
	AlbumArtist string      `json:"albumArtist"`
	Composer    string      `json:"composer"`
	FilePath    string      `json:"filePath"`
}

type TrackNumber struct {
	TrackIndex int `json:"trackIndex"`
	TrackTotal int `json:"trackTotal"`
}

type Library struct {
	AlbumsByTitle map[string]Album `json:"albumsByTitle"`
	TracksByTitle map[string]Track `json:"tracksByTitle"`

	AlbumsById map[int]Album `json:"albumsById"`
	TracksById map[int]Track `json:"tracksById"`
}

func NewLibrary() Library {
	return Library{
		AlbumsByTitle: make(map[string]Album),
		TracksByTitle: make(map[string]Track),
		TracksById:    make(map[int]Track),
		AlbumsById:    make(map[int]Album),
	}
}

type LibraryRequest struct {
	Path string
}

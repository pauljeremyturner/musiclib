package main

type Album struct {
	Id     string     `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
}

type Track struct {
	Id          string         `json:"id"`
	Title       string      `json:"title"`
	Artist      string      `json:"artist"`
	TrackNumber TrackNumber `json:"trackNumber"`
	Album       string      `json:"album"`
	AlbumId     string         `json:"albumId"`
	AlbumArtist string      `json:"albumArtist"`
	Composer    string      `json:"composer"`
	FilePath    string      `json:"filePath"`
}

type TrackNumber struct {
	TrackIndex int `json:"trackIndex"`
	TrackTotal int `json:"trackTotal"`
}
type Library struct {
	Tracks []Track `json:"tracks"`
}

func NewLibrary() Library {
	return Library{
		Tracks:    make([]Track, 0),
	}
}

type LibraryRequest struct {
	Path string
}

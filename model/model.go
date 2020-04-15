package model

type Album struct {
	Id     int
	Title  string
	Artist string
	Tracks []Track
}

type Track struct {
	Id          int
	Title       string
	Artist      string
	TrackNumber TrackNumber
	Album       string
	AlbumArtist string
	Composer    string
	FilePath    string
}

type TrackNumber struct {
	TrackIndex int
	TrackTotal int
}

type Library struct {
	AlbumsByTitle map[string]Album
	TracksByTitle map[string]Track

	AlbumsById map[int]Album
	TracksById map[int]Track
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

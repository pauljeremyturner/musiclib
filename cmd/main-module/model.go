package main

type Album struct {
	Id     int32
	Title  string
	Artist string
	Tracks []Track
}

type Track struct {
	Id          int32
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

	AlbumsById map[int32]Album
	TracksById map[int32]Track
}

func NewLibrary() Library {
	return Library{
		AlbumsByTitle: make(map[string]Album),
		TracksByTitle: make(map[string]Track),
		TracksById:    make(map[int32]Track),
		AlbumsById:    make(map[int32]Album),
	}
}

type LibraryRequest struct {
	Path string
}

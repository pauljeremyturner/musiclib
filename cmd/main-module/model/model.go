package model

type Album struct {
	Title  string
	Artist string
	Tracks []Track
}

type Track struct {
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
	Albums []Album
}

type LibraryRequest struct {
	Path string
}

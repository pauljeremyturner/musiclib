package reader

import (
	"github.com/pauljeremyturner/musiclib/model"
	"gotest.tools/v3/assert"
	"path/filepath"
	"runtime"
	"testing"
)

var (
	_, b, _, _ = runtime.Caller(0)

	// 'testdata'' folder of this project
	testDataDirectory = filepath.Join(filepath.Dir(b), "..", "testdata")
	testTracks        []model.Track
	library           model.Library

	flacTrack model.Track
	mp3Track  model.Track
	oggTrack  model.Track
)

func TestMain(m *testing.M) {
	setupMetaDataTest()

	m.Run()
}

func setupMetaDataTest() {

	md := NewMetaDataReader()

	testTracks = md.ReadMetaData(testDataDirectory)

	titlemap := make(map[string]model.Track)

	for _, t := range testTracks {
		titlemap[t.Title] = t
	}

	flacTrack = titlemap["flac title"]
	mp3Track = titlemap["mp3 title"]
	oggTrack = titlemap["ogg title"]

}

func TestShouldReadArtistMetadataFromFiles(t *testing.T) {

	assert.Equal(t, flacTrack.Artist, "flac artist")
	assert.Equal(t, mp3Track.Artist, "mp3 artist")
	assert.Equal(t, oggTrack.Artist, "ogg artist")

}

func TestShouldReadATitleMetadataFromFiles(t *testing.T) {

	assert.Equal(t, flacTrack.Title, "flac title")
	assert.Equal(t, mp3Track.Title, "mp3 title")
	assert.Equal(t, oggTrack.Title, "ogg title")
}

func TestShouldReadAAlbumMetadataFromFiles(t *testing.T) {

	assert.Equal(t, flacTrack.Album, "album")
	assert.Equal(t, mp3Track.Album, "album")
	assert.Equal(t, oggTrack.Album, "album")
}

func TestShouldReadATrackMetadataFromFiles(t *testing.T) {

	assert.Equal(t, flacTrack.TrackNumber.TrackIndex, 1)
	assert.Equal(t, mp3Track.TrackNumber.TrackIndex, 2)
	assert.Equal(t, oggTrack.TrackNumber.TrackIndex, 3)
}


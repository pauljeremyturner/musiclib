package main

import (
	"path/filepath"
	"runtime"
	"testing"

	"gotest.tools/assert"
)

var (
	_, b, _, _ = runtime.Caller(0)

	// 'testdata'' folder of this project
	TestDataDirectory = filepath.Join(filepath.Dir(b), "..", "testdata")
	testTracks        []Track

	flacTrack Track
	mp3Track  Track
	oggTrack  Track
)

func setupMetaDataTest() {

	md := NewMetaData(TestDataDirectory)
	testTracks = md.ProcessMetadata()

	titlemap := make(map[string]Track)
	for _, track := range testTracks {
		titlemap[track.Title] = track
	}

	flacTrack = titlemap["flac title"]
	mp3Track = titlemap["mp3 title"]
	oggTrack = titlemap["ogg title"]

}

func TestShouldReadArtistMetadataFromFiles(t *testing.T) {

	setupMetaDataTest()

	assert.Equal(t, flacTrack.Artist, "flac artist")
	assert.Equal(t, mp3Track.Artist, "mp3 artist")
	assert.Equal(t, oggTrack.Artist, "ogg artist")

}

func TestShouldReadATitleMetadataFromFiles(t *testing.T) {

	setupMetaDataTest()

	assert.Equal(t, flacTrack.Title, "flac title")
	assert.Equal(t, mp3Track.Title, "mp3 title")
	assert.Equal(t, oggTrack.Title, "ogg title")
}

func TestShouldReadAAlbumMetadataFromFiles(t *testing.T) {

	setupMetaDataTest()

	assert.Equal(t, flacTrack.Album, "album")
	assert.Equal(t, mp3Track.Album, "album")
	assert.Equal(t, oggTrack.Album, "album")
}

func TestShouldReadATrackMetadataFromFiles(t *testing.T) {

	setupMetaDataTest()

	assert.Equal(t, flacTrack.TrackNumber.TrackIndex, 1)
	assert.Equal(t, mp3Track.TrackNumber.TrackIndex, 2)
	assert.Equal(t, oggTrack.TrackNumber.TrackIndex, 3)
}

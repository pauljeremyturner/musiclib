package main

import (
	"gotest.tools/v3/assert"
	"path/filepath"
	"runtime"
	"testing"
)

var (
	_, b, _, _ = runtime.Caller(0)

	// 'testdata'' folder of this project
	testDataDirectory = filepath.Join(filepath.Dir(b), "testdata")
	metaTestTracks    []Track
	metaTestLibrary   Library

	flacTrack Track
	mp3Track  Track
	oggTrack  Track
)

func init() {
	setupMetaDataTest()
}

func setupMetaDataTest() {

	md := NewMetaDataReader()

	metaTestTracks = md.ReadMetaData(testDataDirectory)

	titlemap := make(map[string]Track)

	for _, t := range metaTestTracks {
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

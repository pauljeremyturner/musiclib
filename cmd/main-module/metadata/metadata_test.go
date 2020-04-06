package metadata

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/pauljeremyturner/musiclib/cmd/main-module/model"
	"gotest.tools/assert"
)

func TestShouldReadArtistMetadataFromFiles(t *testing.T) {
	var (
		_, b, _, _ = runtime.Caller(0)

		// TestDataDirectory folder of this project
		TestDataDirectory = filepath.Join(filepath.Dir(b), "..", "testdata")
	)

	tracks := ProcessMetadata(TestDataDirectory)
	titlemap := make(map[string]model.Track)
	for _, track := range tracks {
		titlemap[track.Title] = track
	}

	flacTrack := titlemap["flac title"]
	mp3Track := titlemap["mp3 title"]
	oggTrack := titlemap["ogg title"]

	assert.Equal(t, flacTrack.Artist, "flac artist")
	assert.Equal(t, mp3Track.Artist, "mp3 artist")
	assert.Equal(t, oggTrack.Artist, "ogg artist")

}

func TestShouldReadATrackMetadataFromFiles(t *testing.T) {

}

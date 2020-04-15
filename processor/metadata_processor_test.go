package processor

import (
	"fmt"
	"github.com/pauljeremyturner/musiclib/model"
	"gotest.tools/assert"
	"strconv"
	"testing"
)

var testTracks []model.Track
var mdp MetaDataProcessor
var library model.Library

func getTestTracks() []model.Track {

	var tracks []model.Track

	for i := 0; i < 6; i++ {
		t := model.Track{
			Id:          int(i),
			Title:       "title" + strconv.Itoa(i),
			Artist:      "artist",
			TrackNumber: model.TrackNumber{i, 6},
			Album:       "album",
			AlbumArtist: "album-artist",
			Composer:    "composer",
			FilePath:    fmt.Sprintf("/home/paul/Music/album/track-%2d.mp3", i),
		}
		tracks = append(tracks, t)
	}

	return tracks
}


func setupLibraryTest() {
	testTracks = getTestTracks()
	mdp = NewMetaDataProcessor()
	library = mdp.TransformMetaData(testTracks)
}

func TestMain(m *testing.M) {
	setupLibraryTest()

	m.Run()
}

func TestOrganizeIntoAlbums(t *testing.T) {

	album := library.AlbumsByTitle["album"]

	for i, track := range album.Tracks {
		assert.Equal(t, "title"+strconv.Itoa(i), track.Title)
	}
}

func TestOrganizeByTrackTitle(t *testing.T) {

	for k, v := range library.TracksByTitle {
		assert.Equal(t, k, v.Title)
	}
}

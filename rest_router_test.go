package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gotest.tools/assert"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
)

func init() {

	responseWriter.header = make(map[string][]string)
	mockMetaDataReader := newMockMedaDataReader()

	rr = NewRestRouter(mockMetaDataReader, NewDatabase())

	var request = &http.Request{
		Method: "post",
		Body:   ioutil.NopCloser(bytes.NewReader([]byte("{ \"Path\":\"/home/foo/bar/Music\"}"))),
	}

	rr.LoadLibrary(responseWriter, request)

}

func getTestTrackData() []Track {

	testTracks := []Track{}
	for i := 0; i < 10; i++ {
		t := Track{
			Id:          fmt.Sprintf("TrackTitle%2d", i),
			Title:       fmt.Sprintf("Track Title %2d", i),
			Artist:      "artist",
			TrackNumber: TrackNumber{TrackIndex: i, TrackTotal: 10},
			Album:       "album",
			AlbumArtist: "",
			Composer:    "",
			FilePath:    "",
		}
		testTracks = append(testTracks, t)
	}
	return testTracks
}

func newMockMedaDataReader() MetaDataReader {
	return &mockMetaDataReader{tracks: []Track{}}
}

type mockResponseWriter struct {
	header http.Header
}

type mockMetaDataReader struct {
	tracks []Track
}

func (mrw mockResponseWriter) Header() http.Header {
	return mrw.header
}

var capturedData []byte
var rr RestRouter
var responseWriter mockResponseWriter

func (mrw mockResponseWriter) Write(data []byte) (int, error) {
	capturedData = data[:]
	return len(data), nil
}

func (mmdr mockMetaDataReader) ReadMetaData(path string) []Track {

	testTracks := getTestTrackData
	mmdr.tracks = testTracks()
	return testTracks()
}

func (mrw mockResponseWriter) WriteHeader(code int) {
	//no op
}

func TestGetAlbumByIdRespondsWithCorrectJson(t *testing.T) {

	capturedData = []byte{}
	request := &http.Request{
		Method:     "get",
		RequestURI: "http://localhost/albums/0",
	}

	rr.GetAlbumById(responseWriter, request)

	gotAlbum := &Album{}
	json.Unmarshal(capturedData, gotAlbum)

	log.Println(gotAlbum)
	//add asserts here
}

func TestGetLibraryRespondsWithCorrectJson(t *testing.T) {

	capturedData = []byte{}
	request := &http.Request{
		Method: "get",
	}

	rr.GetLibrary(responseWriter, request)

	gotLibrary := &Library{}

	json.Unmarshal(capturedData, gotLibrary)

	gotTrack := gotLibrary.Tracks[1]

	assert.Equal(t, gotTrack.Title, "Track Title  2")

}

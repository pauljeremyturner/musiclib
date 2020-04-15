package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pauljeremyturner/musiclib/model"
	"github.com/pauljeremyturner/musiclib/processor"
	"github.com/pauljeremyturner/musiclib/reader"
	"gotest.tools/assert"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestMain(m *testing.M) {

	responseWriter.header = make(map[string][]string)
	var metaDataProcessor = processor.NewMetaDataProcessor()
	mockMetaDataReader := newMockMedaDataReader()

	rr = NewRestRouter(mockMetaDataReader, metaDataProcessor)

	var request = &http.Request{
		Method: "post",
		Body:   ioutil.NopCloser(bytes.NewReader([]byte("{ \"Path\":\"/home/foo/bar/Music\"}"))),
	}

	rr.LoadLibrary(responseWriter, request)

	m.Run()
}

func getTestTrackData() []model.Track {

	testTracks := []model.Track{}
	for i := 0; i < 10; i++ {
		t := model.Track{
			Id:          i,
			Title:       fmt.Sprintf("Track Title %2d", i),
			Artist:      "artist",
			TrackNumber: model.TrackNumber{TrackIndex: i, TrackTotal: 10},
			Album:       "album",
			AlbumArtist: "",
			Composer:    "",
			FilePath:    "",
		}
		testTracks = append(testTracks, t)
	}
	return testTracks
}

func newMockMedaDataReader() reader.MetaDataReader {
	return &mockMetaDataReader{tracks: []model.Track{}}
}

type mockResponseWriter struct {
	header       http.Header
}

type mockMetaDataReader struct {
	tracks []model.Track
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

func (mmdr mockMetaDataReader) ReadMetaData(path string) []model.Track {


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

	rr.GetAlbum(responseWriter, request)

	gotAlbum := &model.Album{}
	json.Unmarshal(capturedData, gotAlbum)


	fmt.Println(gotAlbum)
	//add asserts here
}

func TestGetLibraryRespondsWithCorrectJson(t *testing.T) {

	capturedData = []byte{}
	request := &http.Request{
		Method: "get",
	}

	rr.GetLibrary(responseWriter, request)

	gotLibrary := &model.Library{}

	json.Unmarshal(capturedData, gotLibrary)

	gotTrack := gotLibrary.TracksByTitle["Track Title  2"]

	assert.Equal(t, gotTrack.Title, "Track Title  2")

}

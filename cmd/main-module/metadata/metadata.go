package metadata

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/dhowden/tag"
	"github.com/pauljeremyturner/musiclib/cmd/main-module/model"
)

var tracks []model.Track
var trackChan chan model.Track
var wg sync.WaitGroup

func ProcessMetadata(path string) []model.Track {
	trackChan = make(chan model.Track, 1)
	filepath.Walk(path, visit)
	go join()

	wg.Wait()

	close(trackChan)

	fmt.Println(tracks)

	return tracks
}

func visit(path string, info os.FileInfo, err error) error {

	if info.IsDir() && isMusicDir(info.Name()) {
		wg.Add(1)
		go fork(path, trackChan)
	}

	return nil
}

func fork(path string, trackChan chan model.Track) {

	log.Printf("Loading music metadata from %s", path)

	files, _ := ioutil.ReadDir(path)

	defer wg.Done()

	for _, fileinfo := range files {
		filename := filepath.Join(path, fileinfo.Name())

		f, err := os.Open(filename)

		if !isMusicFile(fileinfo.Name()) {
			return
		}

		m, err := tag.ReadFrom(f)
		if err == nil {

			trackIndex, trackTotal := m.Track()

			t := model.Track{
				Title: m.Title(), Album: m.Album(), Artist: m.Artist(),
				TrackNumber: model.TrackNumber{trackIndex, trackTotal},
				AlbumArtist: m.AlbumArtist(), Composer: m.Composer(), FilePath: path,
			}

			trackChan <- t

			f.Close()
		}
	}
}

func join() {
	for track := range trackChan {
		tracks = append(tracks, track)
	}
}

func isMusicFile(path string) bool {

	flac := strings.HasSuffix(path, ".flac")
	mp3 := strings.HasSuffix(path, ".mp3")
	ogg := strings.HasSuffix(path, ".ogg")

	return flac || mp3 || ogg
}

func isMusicDir(name string) bool {

	return !strings.HasPrefix(name, ".")
}

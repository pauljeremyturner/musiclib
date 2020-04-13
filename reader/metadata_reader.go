package reader

import (
	"github.com/dhowden/tag"
	"github.com/pauljeremyturner/musiclib/model"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
)

type metaDataState struct {
	tracks       []model.Track
	trackChan    chan []model.Track
	waitGroupIn  *sync.WaitGroup
	waitGroupOut *sync.WaitGroup
	id           int32
}

type MetaDataReader interface {
	ReadMetaData(path string) []model.Track
}

func NewMetaDataReader() MetaDataReader {
	return &metaDataState{
		tracks:       []model.Track{},
		trackChan:    make(chan []model.Track, 100),
		waitGroupIn:  &sync.WaitGroup{},
		waitGroupOut: &sync.WaitGroup{},
	}
}


func (r *metaDataState) ReadMetaData(path string) []model.Track  {

	_ = filepath.Walk(path, r.visit)

	go r.join()
	r.waitGroupIn.Wait()
	close(r.trackChan)
	r.waitGroupOut.Wait()

	return r.tracks
}

func (r *metaDataState) visit(path string, info os.FileInfo, err error) error {

	if info.IsDir() && isMusicDir(info.Name()) {
		r.waitGroupIn.Add(1)
		r.waitGroupOut.Add(1)

		go r.fork(path)
	} else {
		log.Printf("Not a directory or not a music directory directory: %s", path)
	}

	return nil
}

func (r *metaDataState) fork(path string) {

	log.Printf("Loading music metadata from %s", path)

	files, err := ioutil.ReadDir(path)

	if err != nil {
		log.Printf(err.Error())
		return
	}

	defer r.waitGroupIn.Done()

	ts := []model.Track{}

	for _, fileinfo := range files {
		filename := filepath.Join(path, fileinfo.Name())

		f, err := os.Open(filename)
		if err != nil {
			log.Printf("Could not read track file: error: %s\n", err)
		}

		if !isMusicFile(fileinfo.Name()) {
			log.Printf("Unprocessable file%s\n", fileinfo.Name())
			continue
		}

		m, err := tag.ReadFrom(f)
		if err == nil {

			trackIndex, trackTotal := m.Track()

			t := model.Track{
				Id: atomic.AddInt32(&r.id, 1), Title: m.Title(), Album: m.Album(), Artist: m.Artist(),
				TrackNumber: model.TrackNumber{trackIndex, trackTotal},
				AlbumArtist: m.AlbumArtist(), Composer: m.Composer(), FilePath: filepath.Join(path, filename),
			}

			ts = append(ts, t)
		}

		if f != nil {
			f.Close()
		}
	}
	r.trackChan <- ts

}

func (r *metaDataState) join() {
	for tracks := range r.trackChan {

		for _, t := range tracks {
			r.tracks = append(r.tracks, t)
		}
		r.waitGroupOut.Done()
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

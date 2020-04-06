package main

import (
	"github.com/dhowden/tag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
)

type metaDataState struct {
	tracks       []Track
	trackChan    chan []Track
	waitGroupIn  *sync.WaitGroup
	waitGroupOut *sync.WaitGroup
	path         string
	id           int32
}

type MetaData interface {
	ProcessMetadata() []Track
}

func NewMetaData(inPath string) MetaData {
	return &metaDataState{
		tracks:       []Track{},
		trackChan:    make(chan []Track, 100),
		waitGroupIn:  &sync.WaitGroup{},
		waitGroupOut: &sync.WaitGroup{},
		path:         inPath,
	}
}

func (r *metaDataState) ProcessMetadata() []Track {

	_ = filepath.Walk(r.path, r.visit)

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

	ts := []Track{}

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

			t := Track{
				Id: atomic.AddInt32(&r.id, 1), Title: m.Title(), Album: m.Album(), Artist: m.Artist(),
				TrackNumber: TrackNumber{trackIndex, trackTotal},
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

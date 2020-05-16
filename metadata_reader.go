package main

import (
	"fmt"
	"github.com/dhowden/tag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
)

type metaDataReader struct {
	tracks       []Track
	trackChan    chan []*Track
	waitGroupIn  *sync.WaitGroup
	waitGroupOut *sync.WaitGroup
	id           uint32
	mdb musicDatabase
}
var unknownCounter uint64

type MetaDataReader interface {
	ReadMetaData(path string) []Track
}

func NewMetaDataReader(mdb musicDatabase) MetaDataReader {
	return &metaDataReader{
		tracks:       make([]Track, 0),
		trackChan:    make(chan []*Track, 100),
		waitGroupIn:  &sync.WaitGroup{},
		waitGroupOut: &sync.WaitGroup{},
		mdb: mdb,
	}
}

func (r *metaDataReader) ReadMetaData(path string) []Track {

	_ = filepath.Walk(path, r.visit)

	go r.storeMusicFiles()
	r.waitGroupIn.Wait()
	close(r.trackChan)
	r.waitGroupOut.Wait()

	return r.tracks
}

func (r *metaDataReader) visit(path string, info os.FileInfo, err error) error {

	if info.IsDir() && isMusicDir(info.Name()) {
		r.waitGroupIn.Add(1)
		r.waitGroupOut.Add(1)

		go r.readMusicFilesInDir(path)
	} else {
		log.Printf("Not a directory or not a music directory directory: %s", path)
	}

	return nil
}

func (r *metaDataReader) readMusicFilesInDir(path string) {

	log.Printf("Loading music metadata from %s", path)

	files, err := ioutil.ReadDir(path)

	if err != nil {
		log.Printf(err.Error())
		return
	}

	defer r.waitGroupIn.Done()

	ts := make([]*Track, 0)
	alphanumeric, _ := regexp.Compile("[^a-zA-Z0-9]+")

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


			title := m.Title()
			if len(title) == 0 || title == "" {
				iid := atomic.AddUint64(&unknownCounter, 1)
				title = "Unknown Title" + strconv.FormatUint(iid, 10)
			}
			id := alphanumeric.ReplaceAllString(title, "")
			album := m.Album()
			if len(album) == 0 || album == "" {
				iid := atomic.AddUint64(&unknownCounter, 1)
				album = "Unknown Album" + strconv.FormatUint(iid, 10)
			}
			albumId := alphanumeric.ReplaceAllString(album, "")

			t := &Track{
				Id: id, Title: title, Album: album, AlbumId: albumId,
				Artist: m.Artist(),
				TrackNumber: TrackNumber{trackIndex, trackTotal},
				FilePath: filename,
			}

			ts = append(ts, t)
		}

		if f != nil {
			f.Close()
		}
	}
	r.trackChan <- ts

}

func (r *metaDataReader) storeMusicFiles() {
	for tracks := range r.trackChan {

		uniqueAlbums := make(map[string]*Album)

		for _, t := range tracks {
			_, ok := uniqueAlbums[t.AlbumId]; if !ok {
				uniqueAlbums[t.AlbumId] = &Album{
					Id:     t.AlbumId,
					Title:  t.Album,
					Artist: t.Artist,
				}
			}


		}
		for _, a := range uniqueAlbums {
			fmt.Println("store album", a)
			r.mdb.StoreAlbum(a)
		}


		r.mdb.StoreTracks(tracks)
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

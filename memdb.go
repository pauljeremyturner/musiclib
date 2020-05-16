package main

import (
	"github.com/hashicorp/go-memdb"
)

type musicDatabase struct {
	database *memdb.MemDB
}

func NewDatabase() musicDatabase {

	db, err := memdb.NewMemDB(schema())
	if err != nil {
		panic(err)
	}

	return musicDatabase{
		database: db,
	}
}

func (r musicDatabase) StoreTracks(ts []*Track) {
	txn := r.database.Txn(true)
	for _, t := range ts {
		if err := txn.Insert("track", t); err != nil {
			panic(err)
		}
	}
	txn.Commit()
}


func (r musicDatabase) StoreTrack(t *Track) {
	txn := r.database.Txn(true)
		if err := txn.Insert("track", t); err != nil {
			panic(err)
		}

	txn.Commit()
}


func (r musicDatabase) StoreAlbum(a *Album) {
	txn := r.database.Txn(true)
		if err := txn.Insert("album", a); err != nil {
			panic(err)
		}
	txn.Commit()
}

func (r musicDatabase) LoadAlbumById(albumId string) (*Album, bool) {
	txn := r.database.Txn(false)
	defer txn.Abort()

	a, err := txn.First("album", "id", albumId)
	if a == nil {
		return nil, false
	}
	if err != nil {
		panic(err)
	}

	return a.(*Album), true
}

func (r musicDatabase) LoadTracksByAlbumId(albumId string) ([]*Track, bool) {
	txn := r.database.Txn(false)
	defer txn.Abort()

	it, err := txn.Get("track", "albumId", albumId)
	if err != nil {
		panic(err)
	}
	if it == nil {
		return nil, false
	}
	trs := make([]*Track, 0)
	for tr := it.Next(); tr != nil; tr = it.Next() {
		trs = append(trs, tr.(*Track))

	}
	return trs, true
}

func (r musicDatabase) LoadTracksById(id string) ([]*Track, bool) {
	txn := r.database.Txn(false)
	defer txn.Abort()

	it, err := txn.Get("track", "id", id)
	if err != nil {
		panic(err)
	}
	if it == nil {
		return nil, false
	}
	trs := make([]*Track, 0)
	for tr := it.Next(); tr != nil; tr = it.Next() {
		trs = append(trs, tr.(*Track))

	}
	return trs, true
}

func (r musicDatabase) LoadTracks() ([]*Track, bool) {
	txn := r.database.Txn(false)
	defer txn.Abort()

	it, err := txn.Get("track", "id")
	if err != nil {
		panic(err)
	}
	if it == nil {
		return nil, false
	}
	trs := make([]*Track, 0)
	for tr := it.Next(); tr != nil; tr = it.Next() {
		trs = append(trs, tr.(*Track))

	}
	return trs, true
}

func schema() *memdb.DBSchema {
	return &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			"album": &memdb.TableSchema{
				Name: "album",
				Indexes: map[string]*memdb.IndexSchema{
					"id": &memdb.IndexSchema{
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "Id"},
					},
					"title": &memdb.IndexSchema{
						Name:    "title",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "Title"},
					},
				},
			},
			"track": &memdb.TableSchema{
				Name: "track",
				Indexes: map[string]*memdb.IndexSchema{
					"id": &memdb.IndexSchema{
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "Id"},
					},
					"title": &memdb.IndexSchema{
						Name:    "title",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "Title"},
					},
					"albumId": &memdb.IndexSchema{
						Name:    "albumId",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "AlbumId"},
					},
					"album": &memdb.IndexSchema{
						Name:    "album",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "Album"},
					},
				},
			},
		},
	}
}

# Music Library

Utility to convert a music folder to a json file based on the track metadata of the music files.

The music library can be queried via a ReST api.

This is my first golang project.  I wanted to use what I have learned including 

- interfaces
- unit testing
- encapsulation 
- simple rest service via gorilla mux
- some of the core libraries
- exception handling
- goroutines and channels
- modules
- the go tooling for test covereage
- a dockerised go app
- Hashicorp's go-memdb

It's not finished, I would like to add integration testing via running a docker compose of the app and sending it web requests.  I think the types can be tidied: there probably isn't any need for the metadata processor- it could all go in to the constructor of Library, I probably ignore some errors coming out of I/O calls and surely there is a framework for mocking http requests etc rather than hand coding mock reqtests and output buffers.  


## How does it work?
This is a *fork-join* process.
fork - for each directory, extract the music metadata
join - for all the track metadata objects, combine these into a set of Albums

Anticipating this process to be I/O bound, a set of directories - suspected as each being an album is extracted.
For each directory a goroutine is executed which send the Track metadata objects to a channel.

## How do I use this?

The program expects you music files to be in `~/Music` as this is where I keep my music  on my Fedora system.  You can provide a override directory if your music files are not in this directory.

run the file `main.go`, then call the `librarys` POST method, supplying a folder where your music is located.

## Graphql api

Try this for an example

```
query {
  album(id: "JohnWilliamsPlaysBach") {
    id
    title
    artist
    tracks {
      id
      title
      filePath
    }
  }
}
```

## ReST api

The api uses an in-memory representation of the music library which is stored in a go-memdb database.

These are the resources:

### Library

---

HTTP POST `librarys`
consumes - application/ json
```
{
    "root" : "/home/paul/Music"
}
```

produces - no content

HTTP STATUS: 201 Accepted

---

HTTP GET `librarys`

produces - application/json

example: An excetpt from the json response of complete music library
```
  },
    "Bach : Goldberg Variations, Partitas 5 & 6": {
      "Id": 48,
      "Title": "Bach : Goldberg Variations, Partitas 5 & 6",
      "Artist": "Glenn Gould ",
      "Tracks": [
        {
          "Id": 436,
          "Title": "Goldberg Variations BWV 988 - Aria",
          "Artist": "Glenn Gould",
          "TrackNumber": {
            "TrackIndex": 1,
            "TrackTotal": 18
          },

```


### Albums

---

HTTP GET `/albums/${albumid}`

produces - application/json

path parameters: albumid- The id of the album: assigned during parse of music file metadata

### Tracks

---

HTTP GET `/tracks/${trackid}`

produces - application/json

path parameters: trackid- The id of the album: assigned during parse of music file metadata

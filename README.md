# Music Library

Utility to convert a music folder to a json file based on the track metadata of the music files.

The music library can be queried via a ReST api.

# How does it work?
This is a *fork-join* process.
fork - for each directory, extract the music metadata
join - for all the track metadata objects, combine these into a set of Albums

Anticipating this process to be I/O bound, a set of directories - suspected as each being an album is extracted.
For each directory a goroutine is executed which send the Track metadata objects to a channel.

# How do I run this?

The program expects you music files to be in `~/Music` as is the default on my Fedora system.  You can provide a override directory if your music files are not in this directory.

run the file `musicLibrary.go`, the music metadata json file will be saved in the same directory s where music is located.

# ReST api

The api uses an in-memory representation of the music library.

These are the resources:

### Library

HTTP POST `librarys`
consumes - application/ json
```
{
    "root" : "/home/paul/Music"
}
```

produces - no content

HTTP STATUS: 201 Accepted

### Albums
HTTP GET `/albums/${album name}`

produces - application/json

path parameters: album name- The name of the album

### Artists
HTTP GET `/artists/${artist name}`

produces - application/json

path parameters: artist name- The name of the artist

### Tracks
HTTP GET `/tracks/${track title}`

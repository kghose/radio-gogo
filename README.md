# Radio Go Go

    All we heard is radio ga ga
    Radio goo goo

    - Queen

Simple terminal based radio player that uses [`mpv`][mpv] to play radio streams
whose urls are retrieved from [Radio Browser][radiobrowser]. It is inspired by
[Tera][tera].

[mpv]: https://mpv.io
[radiobrowser]: https://www.radio-browser.info
[tera]: https://github.com/shinokada/tera

# Learning goals

1. Communication between goroutines
1. Sockets (communicate with mpv)
1. Networking basics
1. TUIs

# Learning notes
## The magic of json Decode
 

# User flows



# User interface
1. Basic curses window with a scrolling part and some non-scrolling parts
1. Different "pages" based on what mode we are in

## Modes
### Search
1. Static info:
   1. Currently selected server
   1. Current search criteria
   1. Current station
   1. Current song
   1. mpv status (?)
1. Scrolling info
   1. Recently played history (can add/remove from fav songs)
   1. Stations (when searching)
   1.

### Play
 

## Autoplay
```
1. Select mood
2. Auto play
```

## Search
```
1. Select mood
2. Search by tag and or country
3. Select station from list
```

Auto com

## Station added to blacklist
```
If auto play play a new sttaion
else go back to search page
```

## Favorite station
Add station to current mood.

## Favorite track
Add track details to current mood.

## How moods work
We set a user defined "mood" (Defaults are Happy, Angry, Calm, Working, Writing) 
for the session and attach our saved tracks and stations to this mood. When the 
program is asked to play according to our mood it picks stations based on the 
mood. It can't play individual tracks, but saves the favorite track details in a 
CSV file in case we want to import it to some other application. 

We restrict the number of moods to 5.


## mpv notes

1. Use [JSON IPC](https://mpv.io/manual/master/#json-ipc) for control
1. Start with `mpv --input-ipc-server=/tmp/mpvi.sock --profile=`
1. Play URL: `{ "command": ["loadfile", "URL"], "request_id": 22 }`
1. Pause: `{ "command": ["pause"], "request_id": 22 }`
1. Properties: `{ "command": ["get_property", "PROP"] }`
   path, duration,percent-pos, time-pos,time-remaining, metadata
1. Quit: `{ "command": ["quit"], "request_id": -22 }`



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
 

# User interface

```
[t]ag: ________
[s]tation: _________
|-------------------|
| list item         |
| list item         |
| ...               |
|-------------------|

server: XYZ
```

Typing "t" puts us in the tag entry box. We start typing and the list box lists 
the tags that match the string we are typing. Hitting up/down arrow highlights 
different tags in the list. Hitting enter selects/de-selects list items.

Typing "s" puts us in station list mode. The list changes to show all stations 
matching the selected tag(s). Hitting up/down highlights different stations.  
Hitting enter selects the station and starts playing it. 

On restarting we start with the last state (which is saved under 
`$XDG_CONFIG_HOME/radio-gogo/config.json`) which means what we typed in the tag 
field, selected tags and selected station are saved and restored.


## mpv notes

1. Use [JSON IPC](https://mpv.io/manual/master/#json-ipc) for control
1. Start with `mpv --input-ipc-server=/tmp/mpvi.sock --profile=`
1. Play URL: `{ "command": ["loadfile", "URL"], "request_id": 22 }`
1. Pause: `{ "command": ["pause"], "request_id": 22 }`
1. Properties: `{ "command": ["get_property", "PROP"] }`
   path, duration,percent-pos, time-pos,time-remaining, metadata
1. Quit: `{ "command": ["quit"], "request_id": -22 }`



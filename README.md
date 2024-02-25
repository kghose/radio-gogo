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

# User flows

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

# Learning goals

1. Sockets (communicate with mpv)
2. Networking basics
3. Golang
4. TUIs

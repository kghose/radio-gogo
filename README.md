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

## Struct field tags

# Design decisions: Focus on simplicity

I decided to make a simple app that looks nice (to me).

I initially thought about a tag search, and a station search and so on. But I 
always liked the (original) Google search page the most. Just a simple text box 
where you type in your terms and get results and behind that is a more 
sophisticated search language if you want it.

_I decided to have a single search bar where typing a string does a tag search.  
If you add search directives, it does a more detailed search._

I had initially planned a "responsive" program, where, for example, you could be 
typing a search term and the search results would populate as you typed. This 
would probably be cool, but it brought with it some technical and UX issues.  
What happens if you type a search term, move to the search results and highlight 
something and then more search results come in? Asynchronous events make things 
very complicated.

_I decided search is a blocking action. The UI is not actually blocked: You can 
do a new search while the old search is running (which cancels the old search) 
or quit the app. But it simplifies the UX as well as the code._

# Design decisions: Least dependencies

I did a silly [console game](github.com/kghose/pinman) using nsf/termbox-go and 
liked what I saw there. I especially liked that nsf/termbox-go has few 
dependencies but it would be a lot of work to expand that to a TUI. rivo/tview 
hit the right spot: it is well featured and keeps its deps to within a certain 
circle of what I'd consider well used and stable code. 

nsf/termbox-go is small enough to vendor in, rivo/tview is larger, but feasible 
to vendor in.


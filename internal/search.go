package radio

import (
	"regexp"
	"strings"

	radioBrowser "github.com/kghose/radio-go-go/internal/radio_browser"
)

const searchHelp = `
[yellow]Search[-]

A plain string does a tag search.
e.g. [yellow]classic rock[-] searches for "classic rock"

[yellow]n:[-] searches by name (case insensitive).
[yellow]c:[-] searches by country (case sensitive).
[yellow]t:[-] searches by tag. Tags can be repeated.

e.g. [yellow]n:bbc c:United Kingdom t:pop t:jazz[-]
finds [white]BBC Radio 6 music[-] for us.
`

var re = regexp.MustCompile(`((^|\s)[n|c|t]:)`)

func ParseSearchString(searchStr string) radioBrowser.SearchQuery {
	sq := radioBrowser.SearchQuery{}
	indices := re.FindAllStringIndex(searchStr, -1)
	if len(indices) == 0 {
		sq.TagList = []string{searchStr}
		return sq
	}

	tags := []string{}
	for i, se := range indices {
		key := strings.TrimSpace(searchStr[se[0] : se[1]-1])
		v0 := se[1]
		v1 := len(searchStr)
		if i+1 < len(indices) {
			v1 = indices[i+1][1] - 2
		}
		value := strings.TrimSpace(searchStr[v0:v1])
		switch key {
		case "n":
			sq.Name = value
		case "c":
			sq.Country = value
		case "t":
			tags = append(tags, value)

		}
	}
	if len(tags) > 0 {
		sq.TagList = tags
	}
	return sq
}

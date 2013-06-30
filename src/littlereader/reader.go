package littlereader

import (
	"encoding/json"
	"fmt"
	"github.com/hoisie/web"
	"io/ioutil"
	"os"
)

var folders map[string][]*Source

func readState() {
	state := State{}
	wd, err := os.Getwd()
	fmt.Printf("wd: %s\n", wd)
	file, err := os.Open("state.json")
	if err != nil {
		panic(err)
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &state)
	if err != nil {
		panic(err)
	}
	folders = state.Folders
}

func index() string {
	var s = ""
	var id = 0
	s += `
<!doctype html>
<html>
<head>
<script src="http://code.jquery.com/jquery-1.10.1.min.js"></script>
</head>
<body>
<script>
function hide(s, link) {
	var num = s.split('_')[1];
	$.post('/markAsRead', { href: link });
	$('#' + s).hide();
}
</script>
`
	for folderName, folder := range folders {
		s += fmt.Sprintf("<h2>%s</h2>", folderName)
		for _, source := range folder {
			s += fmt.Sprintf("<h3>%s</h3>", source.Title)
			s += "<ul>"
			for _, entry := range source.Entries {
				if entry.Read {
					continue
				}
				s += fmt.Sprintf(`<li id="entry_%d">`, id)
				s += fmt.Sprintf(`<button onClick="hide('entry_%d', '%s'); return false">Mark Read</button> `, id, entry.Url)
				s += fmt.Sprintf(`<a href="%s">%s</a>`, entry.Url, entry.Title)
				s += "</li>"
				id += 1
			}
			s += "</ul>"
		}
	}
	s += "</body></html>"
	return s
}

func markAsRead(ctx *web.Context) {
	link := ctx.Params["href"]
	fmt.Printf("Marking %s as read\n", link)
	for _, folder := range folders {
		for _, source := range folder {
			for _, entry := range source.Entries {
				if entry.Url == link {
					entry.Read = true
				}
			}
		}
	}
}

func Reader() {
	readState()
	web.Get("/", index)
	web.Post("/markAsRead", markAsRead)
	web.Run("0.0.0.0:9090")
}

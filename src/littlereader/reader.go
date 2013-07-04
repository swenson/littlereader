// Copyright 2013 Christopher Swenson.
//
// This file holds the main web app for little reader.

package littlereader

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hoisie/web"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

var folders map[string][]*Source
var dirty = false

// Read the state from disk.
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

// Generate the main page.
func index() string {
	var buffer bytes.Buffer
	var id = 0
	buffer.WriteString(`
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
`)
	for folderName, folder := range folders {
		buffer.WriteString(fmt.Sprintf("<h2>%s</h2>", folderName))
		for _, source := range folder {
			if !anyNonRead(source) {
				continue
			}
			buffer.WriteString(fmt.Sprintf("<h3>%s</h3>", source.Title))
			buffer.WriteString("<ul>")
			for _, entry := range source.Entries {
				if entry.Read {
					continue
				}
				buffer.WriteString(fmt.Sprintf(`<li id="entry_%d">`, id))
				buffer.WriteString(fmt.Sprintf(`<button onClick="hide('entry_%d', '%s'); return false">Mark Read</button> `, id, entry.Url))
				buffer.WriteString(fmt.Sprintf(`<a href="%s">%s</a>`, entry.Url, entry.Title))
				buffer.WriteString("</li>")
				id += 1
			}
			buffer.WriteString("</ul>")
		}
	}
	buffer.WriteString("</body></html>")
	return buffer.String()
}

// Checks if the source has any unread entries.
func anyNonRead(source *Source) bool {
	for _, entry := range source.Entries {
		if !entry.Read {
			return true
		}
	}
	return false
}

// Handler for marking an entry as read.
func markAsRead(ctx *web.Context) {
	dirty = true
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

// Goroutine for saving the state.
func saver(ticker *time.Ticker) {
	for {
		select {
		case <-ticker.C:
			println("Writing state")
			if dirty {
				dirty = false
				state := State{folders}

				bytes, err := json.Marshal(state)
				if err != nil {
					panic(err)
				}
				err = ioutil.WriteFile("state.json", bytes, 0644)
				if err != nil {
					panic(err)
				}
			}
		}
	}
}

// Goroutine for updating the state.
func updater(ticker *time.Ticker) {
	for {
		select {
		case <-ticker.C:
			println("Updating feeds")

			for _, folder := range folders {
				for _, source := range folder {
					now := time.Now()
					fmt.Printf("Updating feed %s at %s\n", source.Title, source.Url)
					resp, err := http.Get(source.Url)

					if err != nil {
						fmt.Printf("Error: %s\n", err.Error())
						continue
					}
					defer resp.Body.Close()
					body, err := ioutil.ReadAll(resp.Body)
					if err != nil {
						fmt.Printf("Error: %s\n", err.Error())
						continue
					}
					newSource, err := readAtom(now, source.Url, body)
					if err != nil {
						newSource, err = readRss(now, source.Url, body)
						if err != nil {
							fmt.Printf("Could not parse as atom or RSS... skipping\n")
							continue
						}
					}
					updateSource(source, newSource)
				}
			}
			dirty = true
		}
	}
}

// Check for new entries and append them to the existing source.
func updateSource(source *Source, newSource *Source) {
	for _, newEntry := range newSource.Entries {
		var exists = false
		for _, entry := range source.Entries {
			if entry.Url == newEntry.Url {
				exists = true
				break
			}
		}
		if !exists {
			source.Entries = append(source.Entries, newEntry)
		}
	}
}

// Initialize and run the web app.
func Reader() {
	readState()

	saveTicker := time.NewTicker(1 * time.Minute)
	go saver(saveTicker)

	updateTicker := time.NewTicker(12 * time.Hour)
	go updater(updateTicker)

	web.Get("/", index)
	web.Post("/markAsRead", markAsRead)
	web.Run("0.0.0.0:9090")
}

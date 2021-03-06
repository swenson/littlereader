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
	"launchpad.net/goamz/aws"
	"launchpad.net/goamz/s3"
	"net/http"
	"os"
	"sort"
	"sync"
	"time"
)

var folders map[string][]*Source
var dirty = false
var lock = new(sync.Mutex)

const s3_bucket = "swenson_rss"

// Read the state from S3.
func readState() {
	state := State{}
	wd, err := os.Getwd()
	fmt.Printf("wd: %s\n", wd)

	// read from S3
	auth, err := aws.EnvAuth()
	if err != nil {
		panic(err.Error())
	}
	s := s3.New(auth, aws.USEast)
	bucket := s.Bucket(s3_bucket)
	data, err := bucket.Get("rss.json")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Data read from S3\n")

	// old way
	/*
		file, err := os.Open("state.json")
		if err != nil {
			panic(err)
		}
		data, err := ioutil.ReadAll(file)
		if err != nil {
			panic(err)
		}
	*/
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
	var class = 0
	buffer.WriteString(indexTemplate)
	lock.Lock()
	for folderName, folder := range folders {
		buffer.WriteString(fmt.Sprintf("<h2>%s</h2>", folderName))
		for _, source := range folder {
			if !anyNonRead(source) {
				continue
			}
			sort.Sort(source)
			buffer.WriteString(fmt.Sprintf("<h3>%s</h3>", source.Title))
			buffer.WriteString(fmt.Sprintf(`<button onClick="hideAll('source_%d'); return false">Mark all as read</button>`, class))
			buffer.WriteString("<ul>")

			for _, entry := range source.Entries {
				if entry.Read {
					continue
				}
				buffer.WriteString(fmt.Sprintf(`<li id="entry_%d">`, id))
				buffer.WriteString(fmt.Sprintf(`<button class="source_%d" onClick="hide('entry_%d', '%s'); return false">Mark Read</button> `, class, id, entry.Url))
				buffer.WriteString(fmt.Sprintf(`<a href="%s">%s</a>`, entry.Url, entry.Title))
				buffer.WriteString("</li>")
				id += 1
			}
			buffer.WriteString("</ul>")
			class += 1
		}
	}
	lock.Unlock()
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

// Handler for adding a new feed
func addNewFeed(ctx *web.Context) string {
	url := ctx.Params["url"]
	source, err := loadFeed(url)
	if err != nil {
		return err.Error()
	}

	lock.Lock()
	folders["uncategorized"] = append(folders["uncategorized"], source)
	lock.Unlock()

	ctx.Redirect(303, "/")
	return ""
}

// Handler for marking an entry as read.
func markAsRead(ctx *web.Context) {
	link := ctx.Params["href"]
	fmt.Printf("Marking %s as read\n", link)
	lock.Lock()
	dirty = true
	for _, folder := range folders {
		for _, source := range folder {
			for _, entry := range source.Entries {
				if entry.Url == link {
					entry.Read = true
				}
			}
		}
	}
	lock.Unlock()
}

// Goroutine for saving the state.
func saver(ticker *time.Ticker) {
	for {
		select {
		case <-ticker.C:
			println("Writing state")
			lock.Lock()
			if dirty {
				dirty = false
				state := State{folders}

				bytes, err := json.Marshal(state)
				if err != nil {
					panic(err)
				}

				// write to S3
				auth, err := aws.EnvAuth()
				if err != nil {
					panic(err)
				}
				s := s3.New(auth, aws.USEast)
				bucket := s.Bucket(s3_bucket)
				err = bucket.Put("rss.json", bytes, "application/json", s3.ACL("private"))
				if err != nil {
					panic(err)
				}
			}
			lock.Unlock()
		}
	}
}

// Goroutine for updating the state.
func updater(ticker *time.Ticker) {
	for {
		select {
		case <-ticker.C:
			println("Updating feeds")
			lock.Lock()
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
			lock.Unlock()
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

	saveTicker := time.NewTicker(15 * time.Second)
	go saver(saveTicker)

	updateTicker := time.NewTicker(12 * time.Hour)
	go updater(updateTicker)

	web.Get("/", index)
	web.Post("/markAsRead", markAsRead)
	web.Post("/add", addNewFeed)
	web.Run("0.0.0.0:9090")
}

// Copyright 2013 Christopher Swenson.

package littlereader

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

// Our types
type State struct {
	Folders map[string][]*Source
}

// Atom types
type Outline struct {
	Title    string    `xml:"title,attr"`
	Text     string    `xml:"text,attr"`
	Type     string    `xml:"type,attr"`
	XmlUrl   string    `xml:"xmlUrl,attr"`
	HtmlUrl  string    `xml:"htmlUrl,attr"`
	Outlines []Outline `xml:"outline"`
}

type Body struct {
	Outlines []Outline `xml:"outline"`
}
type SubscriptionsXml struct {
	XMLName xml.Name `xml:"opml"`
	Body    Body     `xml:"body"`
}

type Feed struct {
	XMLName xml.Name    `xml:"feed"`
	Title   string      `xml:"title"`
	Id      string      `xml:"id"`
	Updated string      `xml:"updated"`
	Author  string      `xml:"author>name"`
	Entries []AtomEntry `xml:"entry"`
}

type AtomEntry struct {
	Title     string `xml:"title"`
	Link      Link   `xml:"link"`
	Id        string `xml:"id"`
	Published string `xml:"published"`
	Updated   string `xml:"updated"`
	Author    string `xml:"author>name"`
	// SummaryType string `xml:"summary>type,attr"`
	Summary string `xml:"summary"`
	// ContentType string `xml:"content>type,attr"`
	Content string `xml:"content"`
}

type Link struct {
	Href string `xml:"href,attr"`
	Rel  string `xml:"rel,attr"`
}

// RSS types
type Rss struct {
	XMLName  xml.Name  `xml:"rss"`
	Channels []Channel `xml:"channel"`
}

type Channel struct {
	Title         string `xml:"title"`
	Description   string `xml:"description"`
	Link          string `xml:"link"`
	LastBuildDate string `xml:"lastBuildDate"`
	PubDate       string `xml:"pubDate"`
	Ttl           int    `xml:"ttl"`
	Items         []Item `xml:"item"`
}

type Item struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	Link        string `xml:"link"`
	Guid        string `xml:"guid"`
	PubDate     string `xml:"pubDate"`
}

type Source struct {
	LastFetched time.Time
	Title       string
	Url         string
	Folder      string
	Entries     []*Entry
}

type Entry struct {
	Title  string
	Author string
	Url    string
	Read   bool
	Body   string
}

// Import feeds from a Google Reader subscriptions.xml file.
func Import() {
	file, err := os.Open("subscriptions.xml")
	if err != nil {
		panic(err)
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	v := SubscriptionsXml{}
	err = xml.Unmarshal(data, &v)
	if err != nil {
		panic(err)
	}

	subscriptions := flatten(v.Body.Outlines)

	sources := make([]*Source, 0)

	for _, feed := range subscriptions {
		fmt.Printf("Loading feed %s\n", feed)
		now := time.Now()
		resp, err := http.Get(feed)

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
		source, err := readAtom(now, feed, body)
		if err != nil {
			source, err = readRss(now, feed, body)
			if err != nil {
				fmt.Printf("Could not parse as atom or RSS... skipping\n")
				continue
			}
		}
		sources = append(sources, source)
	}
	folders := make(map[string][]*Source)
	folders["uncategorized"] = sources
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

// Import an RSS feed
func readRss(now time.Time, url string, data []byte) (*Source, error) {
	rss := Rss{}
	err := xml.Unmarshal(data, &rss)
	if err != nil {
		return nil, err
	}
	if len(rss.Channels) != 1 {
		return nil, errors.New("RSS does not have exactly 1 channel... skipping")
	}
	source := new(Source)
	source.Folder = "uncategorized"
	source.Url = url
	source.LastFetched = now
	source.Title = rss.Channels[0].Title
	entries := make([]*Entry, 0)
	for _, item := range rss.Channels[0].Items {
		newEntry := new(Entry)
		newEntry.Url = item.Link
		newEntry.Author = ""
		newEntry.Body = ""
		newEntry.Read = false
		newEntry.Title = item.Title
		entries = append(entries, newEntry)
	}
	source.Entries = entries
	return source, nil
}

// Import an Atom feed.
func readAtom(now time.Time, url string, data []byte) (*Source, error) {
	feed := Feed{}
	err := xml.Unmarshal(data, &feed)
	if err != nil {
		return nil, err
	}
	source := new(Source)
	source.Folder = "uncategorized"
	source.Url = url
	source.LastFetched = now
	source.Title = feed.Title
	entries := make([]*Entry, 0)
	for _, entry := range feed.Entries {
		newEntry := new(Entry)
		newEntry.Url = entry.Link.Href
		newEntry.Author = entry.Author
		newEntry.Body = ""
		newEntry.Read = false
		newEntry.Title = entry.Title
		entries = append(entries, newEntry)
	}
	source.Entries = entries
	return source, nil
}

// Convert the slice of outlines to a list of URLs.
func flatten(outlines []Outline) []string {
	ret := make([]string, 0)
	for _, o := range outlines {
		if o.XmlUrl != "" {
			ret = append(ret, o.XmlUrl)
		}
		if len(o.Outlines) > 0 {
			for _, o2 := range flatten(o.Outlines) {
				ret = append(ret, o2)
			}
		}
	}
	return ret
}

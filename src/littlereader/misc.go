package littlereader

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func Reader() {

}

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
	Title       string `xml:"title"`
	Link        Link   `xml:"link"`
	Id          string `xml:"id"`
	Published   string `xml:"published"`
	Updated     string `xml:"updated"`
	Author      string `xml:"author>name"`
	SummaryType string `xml:"summary>type,attr"`
	Summary     string `xml:"summary"`
	ContentType string `xml:"content>type,attr"`
	Content     string `xml:"content"`
}

type Link struct {
	Href string `xml:"href,attr"`
	Rel  string `xml:"rel,attr"`
}

type State struct {
	Folders map[string][]Source
}

type Source struct {
	LastFetched time.Time
	Title       string
	Url         string
	Folder      string
	Entries     []Entry
}

type Entry struct {
	Title  string
	Author string
	Url    string
	Read   bool
	Body   string
}

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

	sources := make([]Source, 0)

	for _, feed := range subscriptions {
		fmt.Printf("Loading feed %s\n", feed)
		resp, err := http.Get(feed)
		now := time.Now()

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
		feed := Feed{}
		err = xml.Unmarshal(body, &feed)
		if err != nil {
			fmt.Printf("Error: %s\n", err.Error())
			continue
		}
		source := Source{}
		source.Folder = "uncategorized"
		source.Url = feed.Id
		source.LastFetched = now
		source.Title = feed.Title
		entries := make([]Entry, 0)
		for _, entry := range feed.Entries {
			newEntry := Entry{}
			newEntry.Url = entry.Link.Href
			newEntry.Author = entry.Author
			newEntry.Body = entry.Content
			newEntry.Read = false
			newEntry.Title = entry.Title
			entries = append(entries, newEntry)
		}
		source.Entries = entries
		sources = append(sources, source)
	}
	folders := make(map[string][]Source)
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

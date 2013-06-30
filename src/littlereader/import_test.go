package littlereader

import (
	"fmt"
	"testing"
	"time"
)

const exampleRss = `
<?xml version="1.0" encoding="UTF-8" ?>
<rss version="2.0">
<channel>
 <title>RSS Title</title>
 <description>This is an example of an RSS feed</description>
 <link>http://www.someexamplerssdomain.com/main.html</link>
 <lastBuildDate>Mon, 06 Sep 2010 00:01:00 +0000 </lastBuildDate>
 <pubDate>Mon, 06 Sep 2009 16:45:00 +0000 </pubDate>
 <ttl>1800</ttl>

 <item>
  <title>Example entry</title>
  <description>Here is some text containing an interesting description.</description>
  <link>http://www.wikipedia.org/</link>
  <guid>unique string per item</guid>
  <pubDate>Mon, 06 Sep 2009 16:45:00 +0000 </pubDate>
 </item>

</channel>
</rss>
`

const exampleAtom = `
<?xml version="1.0" encoding="UTF-8"?>
<feed xmlns="http://www.w3.org/2005/Atom">
  <title>Christopher Swenson</title>
  <id>http://www.caswenson.com/</id>
  <updated>2013-06-01T18:45:38-00:00</updated>
  <author>
    <name>Christopher Swenson</name>
  </author>
<entry><title>Configuration files in Go</title><link href="http://www.caswenson.com/2013_06_01_configuration_files_in_go" rel="alternate" /><id>2013_06_01_configuration_files_in_go</id><published>2013-06-01T18:45:38-00:00</published><updated>2013-06-01T18:45:38-00:00</updated><author><name>Christopher Swenson</name></author><summary type="html">
&lt;p&gt;The other day, I was starting to port an existing service I had into Go.
There were a lot of iss</summary><content type="html">
&lt;p&gt;The other day, I was starting to port an existing service I had into Go.
There were a lot of issues that I had to tackle to get the functionality I wanted,
including being able to run in at least four different environments:
test, dev, stage, and prod.&lt;/p&gt;
</content></entry></feed>`

func TestRss(t *testing.T) {
	_, err := readRss(time.Now(), []byte(exampleRss))
	if err != nil {
		fmt.Println("Could not decode RSS")
		t.FailNow()
	}
}

func TestAtom(t *testing.T) {
	_, err := readAtom(time.Now(), []byte(exampleAtom))
	if err != nil {
		fmt.Println("Could not decode Atom")
		t.FailNow()
	}
}

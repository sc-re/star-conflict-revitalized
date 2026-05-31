package adm

import (
	"encoding/xml"
	"os/exec"
	"strings"
	"time"
)

type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`
	Channel Channel  `xml:"channel"`
}

type Channel struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description CDATA  `xml:"description"`
	Language    string `xml:"language"`
	PubDate     string `xml:"pubDate"`
	Items       []Item `xml:"item"`
}

type Item struct {
	Title         string `xml:"title"`
	PubDate       string `xml:"pubDate"`
	Description   CDATA  `xml:"description"`
	GUID          GUID   `xml:"guid"`
	DisplayInGame int    `xml:"displayInGame"`
}

type GUID struct {
	IsPermaLink string `xml:"isPermaLink,attr"`
	Value       string `xml:",chardata"`
}

type CDATA struct {
	Value string `xml:",cdata"`
}

func importantRssAnnouncement() string {
	out, err := exec.Command("fortune").Output()
	if err != nil {
		return ""
	}
	ret := strings.ReplaceAll(string(out), "\n", "\n<br>")
	return ret
}

func GetRssFeed(fortune bool) (string, error) {
	feed := RSS{
		Version: "2.0",
		Channel: Channel{
			Title:       "Star-Conflict feed",
			Link:        "http://localhost",
			Description: CDATA{""},
			Language:    "en",
			PubDate:     time.Now().UTC().Format(time.RFC1123Z),
			Items:       []Item{},
		},
	}
	if fortune {
		feed.Channel.Items = append(feed.Channel.Items, Item{

			Title:         "Important announcement",
			PubDate:       time.Now().UTC().Format(time.RFC1123Z),
			Description:   CDATA{importantRssAnnouncement()},
			GUID:          GUID{IsPermaLink: "false", Value: "4fcda18177d9f8d926000000"},
			DisplayInGame: 1,
		})
	}

	out, err := xml.MarshalIndent(feed, "", "  ")
	if err != nil {
		return "", err
	}
	ret := xml.Header + string(out)
	return ret, nil
}

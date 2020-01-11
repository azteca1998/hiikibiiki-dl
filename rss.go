package main

type Rss struct {
	Channels []RssChannel `xml:"channel"`
}

type RssChannel struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`

	Image RssImage  `xml:"image"`
	Items []RssItem `xml:"item"`
}

type RssImage struct {
	Title string `xml:"title"`
	Url   string `xml:"url"`
}

type RssItem struct {
	Title     string           `xml:"title"`
	Enclosure RssItemEnclosure `xml:"enclosure"`
}

type RssItemEnclosure struct {
	Url    string `xml:"url,attr"`
	Length int    `xml:"length,attr"`
	Type   string `xml:"type,attr"`
}

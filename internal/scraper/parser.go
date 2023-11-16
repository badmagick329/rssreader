package scraper

import (
    "encoding/xml"
)

type RSSFeed struct {
    Channel struct {
        Title string `xml:"title"`
        Link string `xml:"link"`
        Description string `xml:"description"`
        Language string `xml:"language"`
        Items []RSSItem `xml:"item"`
    } `xml:"channel"`
}

type RSSItem struct {
    Title string `xml:"title"`
    Link string `xml:"link"`
    Description string `xml:"description"`
    PubDate string `xml:"pubDate"`
}

func ParseRSSFeed(data []byte) (*RSSFeed, error) {
    var feed RSSFeed
    err := xml.Unmarshal(data, &feed)
    if err != nil {
        return nil, err
    }
    return &feed, nil
}


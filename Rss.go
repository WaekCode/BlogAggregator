package main

import (
	"context"
	"encoding/xml"
	"html"
	"io"
	"net/http"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}



func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error){

	var rss *RSSFeed 

	req,err1 := http.NewRequestWithContext(
		ctx,
		"GET",
		feedURL,
		nil,
	)


	if err1 != nil{
		return rss,err1

	}

	req.Header.Set(
		"User-Agent",
		"gator",
	)


	client := &http.Client{}

	resp,err2 := client.Do(req)
	if err2 != nil{
		return rss,err2
	}

	defer resp.Body.Close()


	body,err3 := io.ReadAll(resp.Body)
	if err3 != nil{
		return rss,err3
	}

	err4 := xml.Unmarshal(body,&rss)
	if err4 != nil{
		return rss,err4
	}

	
	rss.Channel.Title = html.UnescapeString(rss.Channel.Title)
	rss.Channel.Description = html.UnescapeString(rss.Channel.Description)

	for i,td := range rss.Channel.Item{
		rss.Channel.Item[i].Title = html.UnescapeString(td.Title)
		rss.Channel.Item[i].Description = html.UnescapeString(td.Description)
	}

	return rss,nil

}
package main

import (
	"context"
	"fmt"
)


func scrapeFeeds(s *State) error{
	feed,err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil{
		return err
	} 

	err2 := s.db.MarkFeedFetched(context.Background(),feed.ID)
	if err2 != nil{
		return  err2
	}

	rssfeed,err3 := fetchFeed(context.Background(),feed.Url)
	if err3 != nil{
		return err3
	}

	fmt.Println()
	fmt.Println(">",feed.Name)
	fmt.Println()
	for _,item := range rssfeed.Channel.Item{
		fmt.Println(item.Title)
	}

	return nil

}
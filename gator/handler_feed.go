package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"gator/internal/database"
	"strings"
	"time"

	"github.com/google/uuid"
)

func scrapeFeeds(s *state) error {
	nextFeed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("error getting next feed from db: %v\n", err)
	}

	lastFetchedAt := sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}
	updatedAt := lastFetchedAt.Time
	err = s.db.MarkFeedFetched(context.Background(), database.MarkFeedFetchedParams{ID: nextFeed.ID, LastFetchedAt: lastFetchedAt, UpdatedAt: updatedAt})
	if err != nil {
		return fmt.Errorf("error marking fetched feed: %v\n", err)
	}

	fetchedFeed, err := fetchFeed(context.Background(), nextFeed.Url)
	if err != nil {
		return fmt.Errorf("error getting contents of feed: %v\n", err)
	}
	layout := time.RFC1123Z
	for _, v := range fetchedFeed.Channel.Item {
		pubDate, err := time.Parse(layout, v.PubDate)
		if err != nil {
			fmt.Println("Error parsing publish date. Trying different layout.")
			layout = time.RFC1123
			pubDate, err = time.Parse(layout, v.PubDate)
			if err != nil {
				fmt.Println("Error parsing publish date with alternate layout. Skipping post creation.")
				continue
			}
		}
		if v.Title == "" || v.Link == "" {
			continue
		}
		post, err := s.db.CreatePost(context.Background(),
			database.CreatePostParams{
				ID:          uuid.New(),
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
				Title:       v.Title,
				Url:         v.Link,
				Description: v.Description,
				PublishedAt: pubDate,
				FeedID:      nextFeed.ID,
			})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint \"posts_url_key\"") {
				continue
			}
			fmt.Printf("Unrecognized error adding post to db: %v\n", err)
		} else {
			fmt.Println("=====================================")
			fmt.Printf("Added post to db: %s\n", post.Title)
		}

	}
	return nil
}

func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage: %s <time between reqs>\n", cmd.name)
	}

	duration := cmd.args[0]
	timeBetweenReqs, err := time.ParseDuration(duration)
	if err != nil {
		return fmt.Errorf("unable to determine duration from input: %s\n%v", duration, err)
	}

	fmt.Printf("Collecting feeds every %s\n", duration)
	ticker := time.NewTicker(timeBetweenReqs)

	for ; ; <-ticker.C {
		fmt.Println("Now Scraping feeds...")
		err = scrapeFeeds(s)
		if err != nil {
			ticker.Stop()
			fmt.Printf("unexpected error while scraping feeds: %v\n", err)
		}
	}
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 2 {
		return fmt.Errorf("usage: %s <name of feed> <url>", cmd.name)
	}

	newFeed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
		Url:       cmd.args[1],
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("unable to create feed as current user: %v", err)
	}

	_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    newFeed.ID,
	})
	if err != nil {
		return fmt.Errorf("error following created feed: %v", err)
	}

	fmt.Println("Feed created successfully:")
	printFeed(newFeed)
	fmt.Println()
	fmt.Println("=====================================")

	return nil
}

func handlerFeeds(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return errors.New("too many arguments")
	}

	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't list feeds: %w", err)
	}

	for _, feed := range feeds {
		fmt.Println("-----")
		createdBy, err := s.db.GetUserFromID(context.Background(), feed.UserID)
		if err != nil {
			return fmt.Errorf("error when fetching creator of feed: %w", err)
		}
		fmt.Printf("Name:       %v\n", feed.Name)
		fmt.Printf("URL:        %v\n", feed.Url)
		fmt.Printf("Created By: %s\n", createdBy)
		fmt.Println("-----")
	}

	return nil
}

func printFeed(f database.Feed) {
	fmt.Printf("* ID:            %s\n", f.ID)
	fmt.Printf("* Created:       %v\n", f.CreatedAt)
	fmt.Printf("* Updated:       %v\n", f.UpdatedAt)
	fmt.Printf("* Name:          %s\n", f.Name)
	fmt.Printf("* URL:           %s\n", f.Url)
	fmt.Printf("* UserID:        %s\n", f.UserID)
}

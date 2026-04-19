package main

import (
	"context"
	"fmt"
	"gator/internal/database"
	"strconv"

	"github.com/google/uuid"
)

func handlerBrowse(s *state, cmd command, user database.User) error {
	if len(cmd.args) > 1 {
		return fmt.Errorf("usage: %s <optional request limit integer>\n", cmd.name)
	}
	limit := 2
	if len(cmd.args) == 1 {
		i, err := strconv.Atoi(cmd.args[0])
		if err != nil {
			return fmt.Errorf("limit must be an integer; %v\n", err)
		}
		limit = i
	}

	feeds, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("unable to get followed feeds for user: %w", err)
	}
	feed_ids := make([]uuid.UUID, len(feeds))
	for i, feed := range feeds {
		feed_ids[i] = feed.FeedID
	}

	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		Limit:   int32(limit),
		FeedIds: feed_ids,
	})

	if err != nil {
		return fmt.Errorf("unable to get posts from feeds: %w", err)
	}

	for _, post := range posts {
		fmt.Printf("%s\n", post.Title)
		fmt.Println("-------------")
		fmt.Printf("-- Published at: %v\n", post.PublishedAt)
		fmt.Printf("-- Url: %s\n", post.Url)
		fmt.Println()
	}
	return nil
}

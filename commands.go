package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/WaekCode/gator/internal/config"
	"github.com/WaekCode/gator/internal/database"
	"github.com/google/uuid"
)

type State struct {
	db     *database.Queries
	Config *config.Config
}

type Command struct {
	Name string   // login
	Args []string // <username>
}

type Commands struct {
	c map[string]func(*State, Command) error
}

func (c *Commands) register(name string, f func(*State, Command) error) {

	c.c[name] = f
}

func (c *Commands) run(s *State, cmd Command) error {

	fun := c.c[cmd.Name]
	if fun == nil {
		return fmt.Errorf("Could not run this command")
	}

	err := fun(s, cmd)

	return err

}

func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("no arg was passed")

	}

	_, errs := s.db.GetUserByName(context.Background(), cmd.Args[0])
	if errs != nil {
		return fmt.Errorf("you cant login to a user that doesnt exists...")

	}

	username := cmd.Args[0]

	err := s.Config.SetUser(username)
	if err != nil {
		return err
	}

	fmt.Println("user has been set")
	fmt.Println("current user:", username)

	return nil
}

func HandlerRegister(s *State, cmd Command) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("no arg was passed")
	}

	_, errs := s.db.GetUserByName(context.Background(), cmd.Args[0])
	if errs == nil {
		return fmt.Errorf("user with that name already exists")

	}

	user, err := s.db.CreateUser(
		context.Background(),
		database.CreateUserParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      cmd.Args[0],
		},
	)

	if err != nil {
		return err
	}

	erra := s.Config.SetUser(user.Name)
	if erra != nil {
		return erra
	}

	fmt.Println("user was created")
	fmt.Println("current user:", user.Name)

	fmt.Println(user)
	return nil

}

func HandlerReset(s *State, cmd Command) error {

	err := s.db.ResetUsers(context.Background())
	if err == nil {
		return fmt.Errorf("Users were deleted...")

	}

	fmt.Println("Users were Not deleted")
	return err

}

func HandlerUsers(s *State, cmd Command) error {
	users, err := s.db.ListUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Could not list users")
	}

	if len(users) == 0 {
		return fmt.Errorf("No users found")

	}

	for _, name := range users {
		f := name
		if name == s.Config.CurrentUserName {
			f += " (current)"
			fmt.Println("*", f)

		} else {
			f := name
			fmt.Println("*", f)

		}

	}

	return nil
}

func HandlerAgg(s *State, cmd Command) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("no arg was passed")
	}
	time_between_reqs := cmd.Args[0]

	timeBetweenRequests, err := time.ParseDuration(time_between_reqs)
	if err != nil {
		return err
	}

	fmt.Println("Collecting feeds every", timeBetweenRequests)

	ticker := time.NewTicker(timeBetweenRequests)
	defer ticker.Stop()
	err1 := scrapeFeeds(s) // run immediately
	if err1 != nil {
		fmt.Println("Error scraping feeds:", err1)
	}
	for range ticker.C {
		err := scrapeFeeds(s)
		if err != nil {
			fmt.Println("Error scraping feeds:", err)
		}
	}

	return nil

}

func HandlerAddFeed(s *State, cmd Command, user database.User) error {

	if len(cmd.Args) < 2 {
		return fmt.Errorf("no arg was passed")

	}

	feed, err2 := s.db.CreateFeed(context.Background(),
		database.CreateFeedParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      cmd.Args[0],
			Url:       cmd.Args[1],
			UserID:    user.ID})

	if err2 != nil {
		return err2
	}

	feedfollows, err3 := s.db.CreateFeedFollow(context.Background(),
		database.CreateFeedFollowParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			UserID:    user.ID,
			FeedID:    feed.ID,
		},
	)

	if err3 != nil {
		return err3
	}

	fmt.Println(feedfollows.FeedName)
	fmt.Println(feedfollows.UserName)

	return nil

}

func HandlerFeeds(s *State, cmd Command) error {
	feeds, err := s.db.Listfeeds(context.Background())
	if err != nil {
		return fmt.Errorf("Could not list feeds")

	}

	if len(feeds) == 0 {
		return fmt.Errorf("No feeds found")

	}

	for _, f := range feeds {
		fmt.Println(f.Name)
		fmt.Println(f.Url)
		fmt.Println(f.UserName)
	}

	return nil
}

func HandlerFollow(s *State, cmd Command, user database.User) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("no arg was passed")
	}

	feed, err := s.db.GetFeedFromUrl(context.Background(), cmd.Args[0])

	if err != nil {
		return err
	}

	followerfeed, err3 := s.db.CreateFeedFollow(context.Background(),
		database.CreateFeedFollowParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			UserID:    user.ID,
			FeedID:    feed.ID,
		},
	)

	if err3 != nil {
		return err3
	}

	fmt.Println(followerfeed.FeedName)
	fmt.Println(followerfeed.UserName)

	return nil
}

func HandlerFollowing(s *State, cmd Command, user database.User) error {
	feeds, err2 := s.db.GetFeedFollowsFromUser(context.Background(), user.ID)
	if err2 != nil {
		return err2
	}

	if len(feeds) == 0 {
		return fmt.Errorf("Current user does not follow any feeds")
	}

	for _, f := range feeds {
		fmt.Println(f.FeedName)
	}

	return nil

}

func middlewareLoggedIn(handler func(s *State, cmd Command, user database.User) error) func(*State, Command) error {

	return func(s *State, cmd Command) error {
		if s.Config.CurrentUserName == "" {
			return fmt.Errorf("No users are logged in...")
		}

		user, err := s.db.GetUserByName(context.Background(), s.Config.CurrentUserName)
		if err != nil {
			return err
		}

		return handler(s, cmd, user)

	}

}

func HandlerUnFollow(s *State, cmd Command, user database.User) error {
	feed, err := s.db.GetFeedFromUrl(context.Background(), cmd.Args[0])
	if err != nil {
		return err
	}

	err2 := s.db.UnfollowFeed(context.Background(),
		database.UnfollowFeedParams{
			UserID: user.ID,
			FeedID: feed.ID,
		})

	if err2 != nil {
		return err2
	}

	fmt.Println("Deleting...", feed.Name)

	return nil

}

func HandlerBrowse(s *State, cmd Command, user database.User) error {
	var limit int32
	if len(cmd.Args) < 1 {
		limit = 2
	} else {
		n, err := strconv.Atoi(cmd.Args[0])
		if err != nil {
			return fmt.Errorf("invalid limit: %w", err)
		}
		limit = int32(n)
	}

	posts, err := s.db.GetPostForUser(context.Background(),
		database.GetPostForUserParams{
			UserID: user.ID,
			Limit:  limit,
		},
	)
	if err != nil {
		return err
	}

	for _, p := range posts {
		fmt.Println(p.Title)
		fmt.Println(p.Url)
		fmt.Println(p.Description)
		fmt.Println(p.PublishedAt)
		fmt.Println("-----")
		fmt.Println()
	}

	return nil

}

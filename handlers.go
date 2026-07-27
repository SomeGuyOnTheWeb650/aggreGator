package main

import (
	"context"
	"fmt"
	"time"
	"log"
	"github.com/google/uuid"
	"aggreGator/internal/database"
)

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Arguments) != 1 {
		log.Fatalf("usage: %v <name>", cmd.Name)
	}
	userParam := database.CreateUserParams{
		ID: uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name: cmd.Arguments[0],
	}
	user, err := s.db.CreateUser(context.Background(), userParam)
	if err != nil {
		log.Fatalf("couldn't create user: %w", err)
	}
	err = s.cfg.SetUser(user.Name)
	if err != nil {
		log.Fatalf("couldn't set current user: %w", err)
	}
	
	fmt.Println("User created successfully:")
	printUser(user)
	return nil
}


func handlerLogin(s *state, cmd command) error {
	if len(cmd.Arguments) != 1 {
		log.Fatalf("usage: %s <name>", cmd.Name)
	}
	_, err := s.db.GetUser(context.Background(), cmd.Arguments[0])
	if err != nil {
		log.Fatalf("couldn't find current user: %w", err)
	}
	err = s.cfg.SetUser(cmd.Arguments[0])
	if err != nil {
		log.Fatalf("couldn't set current user: %w", err)
	}
	fmt.Println("User switched successfully!")
	return nil
}

func handlerReset(s *state, cmd command) error {
	if len(cmd.Arguments) != 0 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}
	err := s.db.ClearTable(context.Background())
	if err != nil {
		log.Fatalf("couldn't clear table: %w", err)
	}
	return nil
}

func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		log.Fatalf("error retrieving users: %w", err)
	}
	current := s.cfg.Current_user_name
	for _, item := range users {
		if item.Name == current {
			fmt.Printf(" * %v (current)\n", item.Name)
			continue	
		}
		fmt.Printf(" * %v\n", item.Name)
	}
	return nil
}

func handlerAgg(s *state, cmd command) error {
	location := "http://www.wagslane.dev/index.xml"
	feed, err := fetchFeed(context.Background(), location)
	if err != nil {
		log.Fatalf("error fetching feed: %w", err)
	}
	fmt.Println(feed.Channel.Item)
	return nil
}

func handlerAddFeed(s *state, cmd command) error {
	current_user_name := s.cfg.Current_user_name
	current_user, err := s.db.GetUser(context.Background(), current_user_name)
	if err != nil {
		log.Fatalf("error getting user %w", err)
	}
	feedParams := database.CreateFeedParams{
		ID: uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name: cmd.Arguments[0],
		Url: cmd.Arguments[1],
		UserID: current_user.ID,
	}
	s.db.CreateFeed(context.Background(), feedParams)
	return nil
}

func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		log.Fatalf("error in Getting Feeds: %w", err)
	}
	for _, feed := range feeds {
		fmt.Printf("feed name: %v\n", feed.Name)
		fmt.Printf("feed url: %v\n", feed.Url)
		user, err := s.db.GetUserWithID(context.Background(), feed.UserID)
		if err != nil {
			log.Fatalf("failed getting created by User: %w", err)
		}
		fmt.Printf("feed created by: %v\n", user.Name)
	}
	return nil
}
func printUser(user database.User) {
	fmt.Printf(" * ID:		%v\n", user.ID)
	fmt.Printf(" * Name:	%v\n", user.Name)
}
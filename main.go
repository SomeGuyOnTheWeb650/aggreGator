package main

import _ "github.com/lib/pq"
import (
	"aggreGator/internal/config"
	"aggreGator/internal/database"
	"database/sql"
	"fmt"
	"os"
	"errors"
	"log"
)

type state struct {
	db *database.Queries
	cfg *config.Config
}
type command struct {
	Name string
	Arguments []string
}

type commands struct {
	cmds map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	value, ok := c.cmds[cmd.Name]
	if ok != true {
		return errors.New("command not registered")
	}
	err := value(s, cmd)
	if err != nil {
		return err
	}
	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.cmds[name] = f
}

func main() {
	f, err := os.OpenFile("logs.log", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("something went wrong trying to log: %w", err)
		}
	defer f.Close()
	log.SetOutput(f)
	
	commands := commands{
		cmds: make(map[string]func(*state, command) error),
	}
	
	cfg, err := config.Read()
	if err != nil {
		fmt.Println("Issue with first json read")
		os.Exit(1)
	}
	db, err := sql.Open("postgres", cfg.Db_url)
	if err != nil {
		fmt.Println("error connecting to db: %v", err)
		os.Exit(1)
	}
	defer db.Close()
	dbQueries := database.New(db)

	newState := &state{
		db: dbQueries,
		cfg: &cfg,
	}
	args := os.Args
	if len(args) < 2 {
		fmt.Println("Not enough args provided")
		os.Exit(1)
	}
	
	
	cmd := command{
		Name: args[1],
		Arguments: args[2:],
	}
	
	commands.register("login", handlerLogin)
	commands.register("register", handlerRegister)
	commands.register("reset", handlerReset)
	commands.register("users", handlerUsers)
	commands.register("agg", handlerAgg)
	commands.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	commands.register("feeds", handlerFeeds)
	commands.register("follow", middlewareLoggedIn(handlerFeedsFollow))
	commands.register("following", middlewareLoggedIn(handlerFeedFollowing))
	commands.register("unfollow", middlewareLoggedIn(handlerFeedUnfollow))
	commands.register("browse", middlewareLoggedIn(handlerBrowse))
	err = commands.run(newState, cmd)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	
}

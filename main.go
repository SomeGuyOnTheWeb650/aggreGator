package main

import _ "github.com/lib/pq"
import (
	"aggreGator/internal/config"
	"aggreGator/internal/database"
	"database/sql"
	"github.com/google/uuid"
	"fmt"
	"os"
	"errors"
	"time"
	"context"
	
)

type state struct {
	db *database.Queries
	cfg *config.Config
}
type command struct {
	name string
	arguments []string
}

type commands struct {
	cmds map[string]func(*state, command) error
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.arguments) < 1 {
		return errors.New("didn't receive at least 1 arguments: login")
	}
	
	user, err := s.db.GetUser(context.Background(), cmd.arguments[0])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	
	s.cfg.SetUser(user.Name)
	fmt.Println("User set to", user)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.arguments) < 1 {
		return errors.New("didn't receive at least 1 arguments: register")
	}
	
	current_time := time.Now()
	userParam := database.CreateUserParams{
		ID: uuid.New(),
		CreatedAt: current_time,
		UpdatedAt: current_time,
		Name: cmd.arguments[0],
	}
	user, err := s.db.GetUser(context.Background(), cmd.arguments[0])
	
	if err == nil {
		fmt.Println("user likely exists")
		fmt.Println(err)
		fmt.Println(user.Name)
		os.Exit(1)
	}
	if user.Name != "" {
		fmt.Println("user exists")
		os.Exit(1)
	}
	newUser, err := s.db.CreateUser(context.Background(), userParam)
	

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = handlerLogin(s, cmd)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("User created successfully: ", newUser.Name)
	return nil
}

func (c *commands) run(s *state, cmd command) error {
	value, ok := c.cmds[cmd.name]
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
		fmt.Println("Issue with sql.Open()")
		os.Exit(1)
	}
	dbQueries := database.New(db)

	newState := state{
		db: dbQueries,
		cfg: &cfg,
	}
	args := os.Args
	if len(args) < 2 {
		fmt.Println("Not enough args provided")
		os.Exit(1)
	}
	
	
	cmd := command{
		name: args[1],
		arguments: args[2:],
	}
	
	commands.register("login", handlerLogin)
	commands.register("register", handlerRegister)
	err = commands.run(&newState, cmd)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	
}

package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/thisantm/go-blog-aggregator/internal/database"

	_ "github.com/lib/pq"
	configuration "github.com/thisantm/go-blog-aggregator/internal/config"
)

type state struct {
	db     *database.Queries
	config *configuration.Config
}

type command struct {
	name string
	args []string
}

func main() {
	config, err := configuration.Read()
	if err != nil {
		log.Fatalf("failed to get config: %v\n", err)
	}

	db, err := sql.Open("postgres", config.DBUrl)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	dbQueries := database.New(db)

	currState := state{
		config: &config,
		db:     dbQueries,
	}

	cmds := commands{
		cmdNamesMap: map[string]func(*state, command) error{},
	}

	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("feeds", handlerFeeds)
	cmds.register("follow", middlewareLoggedIn(handlerFollow))
	cmds.register("following", middlewareLoggedIn(handlerFollowing))
	cmds.register("unfollow", middlewareLoggedIn(handlerUnfollow))

	cmd, err := setCommand()
	if err != nil {
		log.Fatalf("invalid command: %v\n", err)
	}

	err = cmds.run(&currState, cmd)
	if err != nil {
		log.Fatalf("command failed: %v\n", err)
	}
}

func setCommand() (command, error) {
	args := os.Args
	if len(args) < 2 {
		return command{}, fmt.Errorf("less than 2 arguments have been passed")
	}

	cmd := command{
		name: args[1],
	}

	if len(args) > 2 {
		cmd.args = args[2:]
	}

	return cmd, nil
}

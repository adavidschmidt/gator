package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/adavidschmidt/blogaggregator/internal/config"
	"github.com/adavidschmidt/blogaggregator/internal/database"
	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Printf("Error geting config: %s\n", err)
		os.Exit(1)
	}
	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		fmt.Printf("Error connecting to database: %s", err)
	}
	dbQueries := database.New(db)

	s := state{
		db:  dbQueries,
		cfg: &cfg,
	}
	cmds := commands{
		handlers: make(map[string]func(*state, command) error),
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
	cmds.register("browse", middlewareLoggedIn(handlerBrowse))

	if len(os.Args) < 2 {
		fmt.Println("Not enough arguments supplied")
		os.Exit(1)
	}

	var cmd command
	cmd.Name = os.Args[1]
	cmd.Args = os.Args[2:]

	err = cmds.run(&s, cmd)
	if err != nil {
		fmt.Printf("error running command %s\n", err)
		os.Exit(1)
	}
}

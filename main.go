package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/WaekCode/gator/internal/config"
	"github.com/WaekCode/gator/internal/database"
	_ "github.com/lib/pq"
)

func main() {

	c, er := config.Read()
	if er != nil {
		fmt.Println(er)
		os.Exit(1)
	}

	dbURL := c.DbURL
	db, errr := sql.Open("postgres", dbURL)
	if errr != nil {
		fmt.Println(errr)
		os.Exit(1)
	}
	defer db.Close()

	dbQueries := database.New(db)

	state := &State{
		db:     dbQueries,
		Config: &c,
	}

	commands := Commands{
		c: make(map[string]func(*State, Command) error),
	}

	commands.register("login", HandlerLogin)
	commands.register("register", HandlerRegister)
	commands.register("reset", HandlerReset)
	commands.register("users", HandlerUsers)
	commands.register("agg", HandlerAgg)
	commands.register("addfeed", middlewareLoggedIn(HandlerAddFeed))
	commands.register("feeds", HandlerFeeds)
	commands.register("follow", middlewareLoggedIn(HandlerFollow))
	commands.register("following", middlewareLoggedIn(HandlerFollowing))
	commands.register("unfollow", middlewareLoggedIn(HandlerUnFollow))
	commands.register("browse", middlewareLoggedIn(HandlerBrowse))

	if len(os.Args) < 2 {
		fmt.Println("Less then two arguments..")
		os.Exit(1)
	}

	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]

	if cmdName == "commands" {
		fmt.Println("Available commands:")
		for k := range commands.c {
			fmt.Println(" - " + k)
		}
		os.Exit(0)
	}

	command := Command{
		Name: cmdName,
		Args: cmdArgs,
	}

	err := commands.run(state, command)
	if err != nil {
		fmt.Println(err)

	}

}

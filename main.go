package main

import (
	"fmt"

	"os"

	"github.com/WaekCode/BlogAggregator/internal/config"
)	


func main() {

	c,er := config.Read()
	if er != nil{
		fmt.Println(er)
	}

	state := &State{
		Config: &c,
	}

	commands := Commands{
		c:  make(map[string]func(*State, Command) error) ,
	}

	commands.register("login",HandlerLogin)
	userInput := os.Args[1:]
	

	if len(os.Args) <= 2{
		fmt.Println("Less then two arguments..")
		os.Exit(1)
	}

	command := Command{
		Name: userInput[0],
		Args: userInput[1:],
	}

	err := commands.run(state,command)
	if err != nil{
		fmt.Println(err)
	}


}
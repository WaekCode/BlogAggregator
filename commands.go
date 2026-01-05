package main

import (
	"errors"
	"fmt"

	"github.com/WaekCode/BlogAggregator/internal/config"
)


type State struct {
    Config *config.Config
}


type Command struct{
	Name string
	Args []string
}

type Commands struct{
	c map[string]func(*State, Command) error

}

func (c *Commands) register(name string, f func(*State, Command) error){
	
	c.c[name] = f
}




func (c *Commands) run(s *State, cmd Command) error {

	fun := c.c[cmd.Name]
	if fun == nil{
		return  fmt.Errorf("Could not run this command")
	}

	fun(s,cmd)

	return  nil

}


func HandlerLogin(s *State, cmd Command) error{
	

	if len(cmd.Args) == 0{
		return  errors.New("The login handler expects a single argument, the username.")
		
	}
	
	username := cmd.Args[0]


	err := s.Config.SetUser(username)
	if err != nil{
		return  err
	}

	fmt.Println("user has been set")

	return nil
}
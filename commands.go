package main

import (
	"github.com/WaekCode/BlogAggregator/internal/config"
	"errors"
)


type state struct {
    Config *config.Config
}


type command struct{
	Name string
	Args []string
}

func handlerLogin(s *state, cmd command) error{
	

	if len(cmd.Args) == 0{
		return  errors.New("The login handler expects a single argument, the username.")
	}
	
	username := cmd.Args[0]



	err := s.Config.SetUser(username)
	if err != nil{
		return  err
	}

	return nil
}
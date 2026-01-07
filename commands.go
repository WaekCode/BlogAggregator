package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"
	"github.com/WaekCode/BlogAggregator/internal/config"
	"github.com/WaekCode/BlogAggregator/internal/database"
	"github.com/google/uuid"
)


type State struct {
	db  *database.Queries

    Config *config.Config
}


type Command struct{
	Name string // login
	Args []string // <username>
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
	
	_, errs := s.db.GetUserByName(context.Background(), cmd.Args[0])
	if errs != nil{
		println("you cant login to a user that doesnt exists...")
		os.Exit(1)
	}


	if len(cmd.Args) == 0{
		return  errors.New("The login handler expects a single argument, the username.")
		
	}
	
	username := cmd.Args[0]


	err := s.Config.SetUser(username)
	if err != nil{
		return  err
	}

	fmt.Println("user has been set")
	fmt.Println("current user:",username)

	return nil
}


func HandlerRegister(s *State, cmd Command) error{
	if len(cmd.Args) == 0{
		return  fmt.Errorf("No name was passed...")
	}	

	_, errs := s.db.GetUserByName(context.Background(), cmd.Args[0])
	if errs == nil{
			fmt.Println("user with that name already exists")
			os.Exit(1)
	}
	


	user,err := s.db.CreateUser(
		context.Background(),
		database.CreateUserParams{
			ID: uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name: cmd.Args[0],

		},
	)
	
	if err != nil{
		return err
	}




	erra := s.Config.SetUser(user.Name)
	if erra != nil{
		return  erra
	}

	fmt.Println("user was created")
	fmt.Println("current user:",s.Config.CurrentUserName)
	
	fmt.Println(user)
	return nil

}
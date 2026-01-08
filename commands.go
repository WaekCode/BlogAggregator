package main

import (
	"context"
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
	if len(cmd.Args) != 1{
		fmt.Println("no arg was passed")
		os.Exit(1)
	}
	
	_, errs := s.db.GetUserByName(context.Background(), cmd.Args[0])
	if errs != nil{
		println("you cant login to a user that doesnt exists...")
		os.Exit(1)
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
	if len(cmd.Args) != 1{
		fmt.Println("no arg was passed")
		os.Exit(1)
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
	fmt.Println("current user:",user.Name)
	
	fmt.Println(user)
	return nil

}


func HandlerReset(s *State, cmd Command) error{
	
	err := s.db.ResetUsers(context.Background())
	if err == nil{
		fmt.Println("Users were deleted...")
		return err
	}

	fmt.Println("Users were Not deleted")
	return err

}


func HandlerUsers(s *State, cmd Command) error{
	users,err := s.db.ListUsers(context.Background())
	if err != nil{
		fmt.Println("Could not list users")
		return err
	}

	if len(users) == 0{
		fmt.Println("No users found")
		return nil
	}


	for _,name := range users{
		f := name
		if name == s.Config.CurrentUserName{
			f += " (current)"
			fmt.Println("*",f)

			}else{
			f := name
			fmt.Println("*",f)
		
		}

		
	}

	return nil
}


func HandlerAgg(s *State, cmd Command) error{

	rss,err := fetchFeed(context.Background(),"https://www.wagslane.dev/index.xml")
	if err != nil{
		return  err
	}

	fmt.Println(rss)
	
	return nil

}
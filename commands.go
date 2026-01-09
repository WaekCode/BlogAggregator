package main

import (
	"context"
	"fmt"
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

	err := fun(s,cmd)

	return  err

}


func HandlerLogin(s *State, cmd Command) error{
	if len(cmd.Args) < 1{
		return fmt.Errorf("no arg was passed")
	
	}
	
	_, errs := s.db.GetUserByName(context.Background(), cmd.Args[0])
	if errs != nil{
		return fmt.Errorf("you cant login to a user that doesnt exists...")
		
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
	if len(cmd.Args) < 1{
		return fmt.Errorf("no arg was passed")
	}

	_, errs := s.db.GetUserByName(context.Background(), cmd.Args[0])
	if errs == nil{
			return fmt.Errorf("user with that name already exists")
		
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
		return fmt.Errorf("Users were deleted...")

	}

	fmt.Println("Users were Not deleted")
	return err

}


func HandlerUsers(s *State, cmd Command) error{
	users,err := s.db.ListUsers(context.Background())
	if err != nil{
		return fmt.Errorf("Could not list users")
	}

	if len(users) == 0{
		return fmt.Errorf("No users found")
		
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


func HandlerAddFeed(s *State, cmd Command) error{
	
	if len(cmd.Args) < 2{
		return fmt.Errorf("no arg was passed")
	
	}
	curUser := s.Config.CurrentUserName
	if curUser == ""{
		return fmt.Errorf("no user is logged in")
	}else{
	user,err := s.db.GetUserByName(context.Background(),curUser)

	feed,err := s.db.CreateFeed(context.Background(),
								database.CreateFeedParams{
								ID: uuid.New(),
								CreatedAt: time.Now(),
								UpdatedAt: time.Now(),
								Name: cmd.Args[0],
								Url: cmd.Args[1],
								UserID: user.ID,})	

	if err != nil{
		return err
	}

	fmt.Println(feed)
	}

	return  nil							

}




func HandlerFeeds(s *State, cmd Command) error{
	feeds,err := s.db.Listfeeds(context.Background())
	if err != nil{
		return fmt.Errorf("Could not list feeds")
		
	}

	if len(feeds) == 0{
		return fmt.Errorf("No feeds found")
	
	}

	for _,f := range feeds{
		fmt.Println(f.Name)
		fmt.Println(f.Url)
		fmt.Println(f.UserName)
	}

	return nil
}

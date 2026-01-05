package config 

import (
	"os"
	"path/filepath"
	"encoding/json"
)

type Config struct {
	DbURL string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`

}

func getConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	configPath := filepath.Join(homeDir, ".gatorconfig.json")
	return configPath, nil
}

func Read() (Config,error){
	path,ers := getConfigPath()
	if ers != nil{
		return Config{},ers
	}

	res,err := os.ReadFile(path)
	if err != nil{
		return Config{},err
	}

	var c Config
	err2 := json.Unmarshal(res,&c)

	if err != nil{
		return Config{},err2
	}

	return  c,nil


}

func (c Config) SetUser(username string) error{
	c.CurrentUserName = username

	path,err2 := getConfigPath()
	if err2 != nil{
		return err2
	}

	jsonData,err := json.Marshal(c)
	if err != nil{
		return err
	}

	er := os.WriteFile(path,jsonData,0644)

	if er != nil{
		return  er
	}

	return nil

}


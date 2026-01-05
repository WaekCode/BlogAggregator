package main

import (
	"github.com/WaekCode/BlogAggregator/internal/config"
	"fmt"


)	


func main() {

	c,er := config.Read()
	if er != nil{
		fmt.Println(er)
	}

	c.SetUser("Weak")

	c2,er2 := config.Read()
	if er2 != nil{
		fmt.Println(er2)
	}

	fmt.Println(c2)

}
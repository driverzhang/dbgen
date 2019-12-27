package main

import (
	"fmt"
	"github.com/driverzhang/go-mongo-crud-template/crud"
	"os"
)

func main() {
	args := os.Args
	if len(args) == 1 {
		fmt.Println("Bygo error args")
		return
	}

	command := args[1]
	switch command {
	case "db2mongo":
		_, err := crud.GetMongoCrudTemplate()
		if err != nil {
			fmt.Println("bygo", err.Error())
		} else {
			fmt.Println("bygo", "gen mongo crud template struct from clipboard success")
		}
	}

}

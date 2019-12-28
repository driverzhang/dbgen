package main

import (
	"fmt"
	"github.com/driverzhang/dbgen/crud"
	"os"
)

func main() {
	args := os.Args
	if len(args) == 1 {
		fmt.Println("dbgen error args")
		return
	}

	command := args[1]
	switch command {
	case "db2mongo":
		_, err := crud.GetMongoCrudTemplate()
		if err != nil {
			fmt.Println("dbgen", err.Error())
		} else {
			fmt.Println("dbgen", "gen mongo crud template struct from clipboard success")
		}
	}

}

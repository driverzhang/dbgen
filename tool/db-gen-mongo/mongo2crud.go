package db_gen_mongo

import (
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/urfave/cli"
	"strings"
)

func Mongo2Crud(c *cli.Context) (err error) {
	j, err := clipboard.ReadAll()
	if err != nil {
		return
	}
	
	r, err := mongo2Crud(j)
	if err != nil {
		return
	}
	err = clipboard.WriteAll(r)
	fmt.Println("gen mongo db-gen-mongo template struct from clipboard success")
	return
}

func mongo2Crud(input string) (rsp string, err error) {
	firstWord := ""
	for _,v:=range []byte(input) {
		firstWord = strings.ToLower(string(v))
		break
	}
	
	
	t := &TableOptions{
		N: firstWord,
		Name: input,
	}
	t.setName(strings.ToLower(input))
	rsp ,err = t.getMongoCrudTemplate()
	return
}

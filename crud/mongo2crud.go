package crud

import "github.com/atotto/clipboard"

func Mongo2Crud() (err error) {
	j, err := clipboard.ReadAll()
	if err != nil {
		return
	}

	r, err := mongo2Crud(j)
	if err != nil {
		return
	}
	err = clipboard.WriteAll(r)
	return
}

func mongo2Crud(input string) (rsp string, err error) {

	return
}

package main

import (
	"context"
	"fmt"

	"github.com/forgetaboutitapp/forget-about-it/server"
	"github.com/forgetaboutitapp/forget-about-it/server/pkg/dbUtils"
	"github.com/forgetaboutitapp/forget-about-it/server/pkg/sql_queries"
)

func main() {
	fmt.Println(server.DBFilename)
	db, err := dbUtils.OpenDatabase(context.Background())
	if err != nil {
		panic(err)
	}
	q := sql_queries.New(db)

	adminUser, err := q.GetUser(context.Background(), 0)
	if err != nil {
		panic(err)
	}
	if len(adminUser) > 1 {
		panic(fmt.Sprintf("There should not be more than 1 user but there were %d", len(adminUser)))
	} else if len(adminUser) == 0 {
		mnemonic, err := dbUtils.AddUser(q)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Your 12 word mnemonic is:\n%s", mnemonic)
	} else if len(adminUser) == 1 {
		mnemonic, err := dbUtils.AddLogin(context.TODO(), q, adminUser[0])
		if err != nil {
			panic(err)
		}
		fmt.Printf("Your 12 word mnemonic is:\n%s", mnemonic)
	}
}

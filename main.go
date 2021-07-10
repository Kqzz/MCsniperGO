package main

import (
	"fmt"

	"github.com/kqzz/mcgo"
)

func main() {
	var accounts []mcgo.MCaccount
	accStrs, err := readLines("accounts.txt")
	if err != nil {
		panic(err)
	}
	
	accounts = loadAccSlice(accStrs)
	fmt.Println(accounts)
}

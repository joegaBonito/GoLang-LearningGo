package main

import (
	"fmt"

	"github.com/joegabonito/learngo/bankdictionaryproject/accounts"
)

func main() {
	account := accounts.NewAccount("Nico")
	account.Deposit(10)
	fmt.Println(account.Balance())
	err := account.Withdraw(20)
	if err != nil {
		fmt.Println(err)
	}
	account.ChangeOwner("Lesley")
	fmt.Println(account.Balance(), account.GetOwner())
	fmt.Println(account)
}

package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/CHIKUWAODEN/monkey-for-c95/repl"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello %s! This is Monkey programming language!\n",
		user.Username)
	fmt.Printf("Feel free to type in commands.\n")
	repl.Start(os.Stdin, os.Stdout)
}

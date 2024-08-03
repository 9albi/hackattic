package main

import "github.com/9albi/hackattic/cmd"

func main() {
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}

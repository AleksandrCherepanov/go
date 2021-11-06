package main

import (
	"fmt"

	"alexanderlab.club/game_priority/steam"
)

func main() {
	// response, err := steam.GetUserGame("*")
	// response, err := steam.GetUserWishlist("*")
	response, err := steam.GetGameDetails("*")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(response))
}

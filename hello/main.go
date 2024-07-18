package main

import (
	"fmt"

	"karasev.ru/greetings"
)

func main() {
	message := greetings.Hello("Gladys")
	fmt.Println(message)
}

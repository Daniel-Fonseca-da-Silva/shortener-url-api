package main

import "fmt"

func main() {
	fmt.Println(SayMsg("Hello"))
}

func SayMsg(msg string) (say string) {
	say = msg
	return
}

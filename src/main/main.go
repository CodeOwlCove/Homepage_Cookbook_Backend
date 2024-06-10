package main

import "fmt"

func main() {

	fmt.Println("Starting CookBook...")

	initDBDriver()

	initCookBookController()
}

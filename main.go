package main

import ("fmt"
	"strings")

func main() {
fmt.Println("Hello, World!")
}

func cleanInput(text string) []string {
	result := strings.Split(strings.Trim(strings.ToLower(text)," ")," ")
	return result
}

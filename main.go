package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	userInput := getUserInput(reader)
	fmt.Println(userInput)
	handleInput()
	fmt.Println("Exiting...")
}

func getUserInput(r *bufio.Reader) string {
	line, err := r.ReadString('\n')
	if err != nil {
		panic(err)
	}
	return strings.ToLower(strings.TrimSpace(line))
}

func handleInput() {
	reader := bufio.NewReader(os.Stdin)
	action := ""
	for action != "exit" {

		action = getUserInput(reader)
		switch action {
		case "help":
		case "start":
			getRandomArticle()
		case "exit":
		default:
			fmt.Println("Invalid action")
		}
	}

}
func getRandomArticle() {
	url := "https://en.wikipedia.org/w/api.php?"
	var params = map[string]string{
		"action":       "query",
		"format":       "json",
		"generator":    "random",
		"grnnamespace": "0",
		"prop":         "links",
		"pllimit":      "max",
		"plnamespace":  "0",
	}

	for k, v := range params {
		fmt.Printf("Key %s has value %s\n", k, v)
		var param = fmt.Sprintf("&%s=%s", k, v)
		url += param
	}
	fmt.Println(url)
}

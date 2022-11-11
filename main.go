package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
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
func getRandomArticle() Response {
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
	res, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	byteSlice, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	r := processResponse(byteSlice)
	return r
}
func processResponse(b []byte) Response {
	var res Response
	err := json.Unmarshal(b, &res)
	if err != nil {
		log.Println(err)
	}
	return res
}

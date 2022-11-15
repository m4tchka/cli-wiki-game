package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func main() {
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
			res := getRandomArticle()
			var rA Page
			for _, v := range res.Query.Pages {
				rA = v
			}
			displayLinks(rA)
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
		// fmt.Printf("Key %s has value %s\n", k, v)
		var param = fmt.Sprintf("&%s=%s", k, v)
		url += param
	}
	fmt.Println("URL:", url)
	res, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	// fmt.Println("Response:", res)
	byteSlice, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	r := processResponse(byteSlice)
	// fmt.Println("ACTUAL RESPONSE -------->", r)
	return r
}
func processResponse(b []byte) Response {
	var res Response
	err := json.Unmarshal(b, &res)
	if err != nil {
		panic(err)
	}
	// fmt.Println("res: ", res)
	return res
}

func getSpecificArticle() Response {
	url := "https://en.wikipedia.org/w/api.php?"
	var params = map[string]string{
		"action": "query",
		"format": "json",
		// "generator":    "random",
		// "grnnamespace": "0",
		"prop":        "links",
		"pllimit":     "max",
		"plnamespace": "0",
	}

	for k, v := range params {
		// fmt.Printf("Key %s has value %s\n", k, v)
		var param = fmt.Sprintf("&%s=%s", k, v)
		url += param
	}
	fmt.Println("URL:", url)
	res, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	// fmt.Println("Response:", res)
	byteSlice, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	r := processResponse(byteSlice)
	// fmt.Println("ACTUAL RESPONSE -------->", r)
	return r
}
func displayLinks(p Page) {
	fmt.Printf("\nLinks from the page `%s`:\n", p.Title)
	for _, v := range p.Links {
		fmt.Println(v.Title)
	}
}

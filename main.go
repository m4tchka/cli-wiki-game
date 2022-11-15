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

func getUserInput() string {
	r := bufio.NewReader(os.Stdin)
	line, err := r.ReadString('\n')
	if err != nil {
		panic(err)
	}
	return strings.ToLower(strings.TrimSpace(line))
}

func handleInput() {
	action := ""
	for action != "exit" {

		action = getUserInput()
		switch action {
		case "help":
		case "start":
			res := getRandomArticle()
			var rA Page
			for _, v := range res.Query.Pages {
				rA = v
			}
			links := getAndDisplayLinks(rA)
			choice := getUserInput()
			isLinkValid := checkIsLinkValid(links, choice)
			if isLinkValid {
				getSpecificArticle(choice)
			} else {
				"reprompt"
			}

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
		var param = fmt.Sprintf("&%s=%s", k, v)
		url += param
	}
	fmt.Println("URL:", url)
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
		panic(err)
	}
	// fmt.Println("res: ", res)
	return res
}
func getSpecificArticle(t string) Response {
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
func checkIsLinkValid(linkSlice []Link, choice string) bool {
	for _, v := range linkSlice {
		if strings.ToLower(v.Title) == choice {
			return true
		}
	}
	return false
}
func getAndDisplayLinks(p Page) []Link {
	fmt.Printf("\nLinks from the page `%s`:\n", p.Title)
	for _, v := range p.Links {
		fmt.Println(v.Title)
	}
	return p.Links
}

/*
 get a list of valid links (and display them)
 get user input (of a link that they want to "click" on)
 check that link is valid
 fetch that link
 get a list of valid links from the fetched link...
*/

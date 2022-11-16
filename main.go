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

func getUserInput() string { // Prompt for option, return user input as string
	r := bufio.NewReader(os.Stdin)
	line, err := r.ReadString('\n')
	if err != nil {
		panic(err)
	}
	return strings.ToLower(strings.TrimSpace(line))
}

func handleInput() { // Function to handle user input and call corresponding functions

}
func getRandomArticle() Response { // Function to get a random article and return a response struct
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
func processResponse(b []byte) Response { // Function to turn a byte slice into a response
	var res Response
	err := json.Unmarshal(b, &res)
	if err != nil {
		panic(err)
	}
	// fmt.Println("res: ", res)
	return res
}
func getSpecificArticle(t string) Response { // FIXME: Function to get a specific article (based on user input) and return a response struct
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
func checkIsLinkValid(linkSlice []Link, choice string) bool { // Function to check that the user has inputted a valid link (based on the current page)
	for _, v := range linkSlice {
		if strings.ToLower(v.Title) == choice {
			return true
		}
	}
	return false
}
func getAndDisplayLinks(p Page) []Link { // Function that prints all the links in a page to the console, in a easier-to-read format.
	fmt.Printf("\nLinks from the page `%s`:\n", p.Title)
	for _, v := range p.Links {
		fmt.Println(v.Title)
	}
	return p.Links
}
func getPageFromResponse(res Response) Page {
	var p Page
	for _, v := range res.Query.Pages {
		p = v
	}
	return p
}

/*
upon typing getRandom- get a random article - set it's title to be the target page
then, upon typing start, get another random article (the start page)

 get a list of links from the current page (and display them)
 get user input (of a link that they want to "click" on)
 check that link is valid
 fetch that link
 set that link to the current page
 get a list of valid links from the current page...


 ...
 check that link is valid
 if that links is also the target page, end the game and print the count, (and time.Now?)
*/

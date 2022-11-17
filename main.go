package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
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
	return strings.TrimSpace(line)
}

func handleInput() { // Function to handle user input and call corresponding functions
	var action string
	var target string
	var startTime time.Time
	fmt.Println("Welcome to wiki game CLI !")
	instructions := fmt.Sprintf("%-20s- Get a random target page, or if there already is a target, rerolls for a new random target\n%-20s- With a target, get a random start page\n%-20s- Display available commands\n%-20s- Exit the game\n", "get random target", "start", "help", "exit")
	fmt.Println(instructions)
	for action != "exit" {
		var startPage Page
		count := 0
		action = strings.ToLower(getUserInput())
		switch action {
		case "get random target":
			targetArticleResponse := getRandomArticleWithLinks(false)
			targetPage := getPageFromResponse(targetArticleResponse)
			target = targetPage.Title
			fmt.Printf("\nTarget page:        %s\nTarget description: %v\n", target, targetPage.Extract)
		case "get specific target":
			fmt.Println("Not implemented yet")
		case "start":
			if target != "" {
				startArticleResponse := getRandomArticleWithLinks(true)
				startPage = getPageFromResponse(startArticleResponse)
				startTime = time.Now()
				count = startGame(startPage, target) // This when the actual gameplay loop starts
				totalDuration := time.Since(startTime)
				stringDuration := totalDuration.String()
				fmt.Println(">>> Congratulations ! <<<")
				fmt.Printf("You took %d jumps and %s to get to the page %s from %s\n", count, stringDuration, target, startPage.Title)
			} else { //TODO: Add a "get target manually" function to set a specific target page
				fmt.Println("Please set a target page first, with 'get random target' or 'get specific target'!")
			}
		case "help":
			fmt.Println(instructions)
		case "exit":
			action = "exit"
		default:
			fmt.Println("Invalid action !\nType 'help' for available commands.")
			// fmt.Println("Invalid action !\n'get target' - Get a random target page\n'start' - With a target, gets a random start page\n'exit' - exits the game.")
		}

	}
}
func getRandomArticleWithLinks(b bool) Response { // Function to get a random article and return a response struct
	baseURL := "https://en.wikipedia.org/w/api.php?"
	var params = map[string]string{ // Default params - get links
		"action":       "query",
		"format":       "json",
		"generator":    "random",
		"grnnamespace": "0",
		"prop":         "links",
		"pllimit":      "max",
		"plnamespace":  "0",
	} // https://en.wikipedia.org/w/api.php?&generator=random&grnnamespace=0&prop=links&pllimit=max&plnamespace=0&action=query&format=json
	if !b { // If function is passed false, change params to extracts
		params["prop"] = "extracts"
		params["exintro"] = ""
		params["explaintext"] = ""
		params["redirects"] = "1"
		delete(params, "pllimit")
		delete(params, "plnamespace")
	} // https://en.wikipedia.org/w/api.php?&exintro=&action=query&grnnamespace=0&prop=extracts&explaintext=&redirects=1&format=json&generator=random
	for k, v := range params {
		var param = fmt.Sprintf("&%s=%s", k, v)
		baseURL += param
	}
	fmt.Println("random article URL:", baseURL)
	res, err := http.Get(baseURL)
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
	// TODO: WARNING: Does not work with spaces in the name - %20 isntead ?
	fmt.Println("GetSpecificArticle called -----------------")
	baseURL := "https://en.wikipedia.org/w/api.php?"
	var params = map[string]string{
		"action": "query",
		"format": "json",
		// "generator":    "random",
		// "grnnamespace": "0",
		"prop":        "links",
		"pllimit":     "max",
		"plnamespace": "0",
		"titles":      url.QueryEscape(t),
	}

	for k, v := range params {
		// fmt.Printf("Key %s has value %s\n", k, v)
		var param = fmt.Sprintf("&%s=%s", k, v)
		baseURL += param
	}
	fmt.Println("specific article URL:", baseURL)
	res, err := http.Get(baseURL)
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
func checkIsLinkValid(linkSlice []Link, choice string) (bool, string) { // Function to check that the user has inputted a valid link (based on the current page)
	for _, v := range linkSlice {
		if strings.EqualFold(v.Title, choice) {
			return true, v.Title
		} // WARNING: SUS
	}
	return false, ""
}
func getAndDisplayLinks(p Page) []Link { // Function that prints all the links in a page to the console, in a easier-to-read format.
	fmt.Printf("\nLinks from the page `%s`:\n", p.Title)
	for i, v := range p.Links {
		fmt.Println(i, v.Title)
	}
	return p.Links
}
func getPageFromResponse(res Response) Page { // Function to extract the page from the response
	var p Page
	for _, v := range res.Query.Pages { // Loop is needed since pages are returned in a map of ints to Pages, but the int key is not known beforehand. Only a single page is returned.
		p = v
	}
	return p
}
func startGame(p Page, t string) int {
	count := 0
	currentPage := p
	for currentPage.Title != t {
		currentLinks := getAndDisplayLinks(currentPage)
		fmt.Printf("Current page:   %s\n, Current target: %s\n", currentPage.Title, t)
		userChoice := getUserInput()
		isLinkValid, validChoice := checkIsLinkValid(currentLinks, userChoice)
		if isLinkValid {
			selectedArticle := getSpecificArticle(validChoice)
			currentPage = getPageFromResponse(selectedArticle)
			count++
		} else {
			fmt.Println("Link not found !")
		}
	}
	return count
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

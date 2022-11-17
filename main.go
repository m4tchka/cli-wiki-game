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
	instructions := fmt.Sprintf("%-25s- Get a random target page, or if there already is a target, rerolls for a new random target\n%-25s- Get a specified target page, by entering its name on the line below\n%-25s- With a target, get a random start page\n%-25s- Display available commands\n%-25s- Exit the game\n", "get random target", "get specific target", "start", "help", "exit")
	fmt.Println(instructions)
	for action != "exit" {
		var startPage Page
		var count int
		var linksTaken []string
		action = strings.ToLower(getUserInput())
		switch action {
		case "get random target":
			targetArticleResponse := getRandomArticleWithLinks(false)
			targetPage := getPageFromResponse(targetArticleResponse)
			target = targetPage.Title
			fmt.Printf("\nTarget page:\n%s\n%s\n", target, targetPage.Extract)
		case "get specific target":
			fmt.Print("Define a target page to navigate to: ")
			targetArticleResponse := getSpecificArticle(getUserInput())
			targetPage := getPageFromResponse(targetArticleResponse)
			target = targetPage.Title
			fmt.Printf("\nTarget page:\n%s\n%s\n", target, targetPage.Extract)
		case "start":
			if target != "" {
				startArticleResponse := getRandomArticleWithLinks(true)
				startPage = getPageFromResponse(startArticleResponse)
				fmt.Printf("Starting page:\n%s\n%s\n", startPage.Title, startPage.Extract)
				startTime = time.Now()
				count, linksTaken = startGame(startArticleResponse, target) // This when the actual gameplay loop starts
				totalDuration := time.Since(startTime)
				stringDuration := totalDuration.String()
				fmt.Println("\n>>> Congratulations ! <<<")
				fmt.Printf("\nYou took %d jumps and %s to get to the page %s from %s\n", count, stringDuration, target, startPage.Title)
				fmt.Printf("\nYou took the route: %s\n", strings.Join(linksTaken, " -> "))
			} else {
				fmt.Println("Please set a target page first, with 'get random target' or 'get specific target'!")
			}
		case "help":
			fmt.Println(instructions)
		case "exit":
			action = "exit"
		default:
			fmt.Println("Invalid action !\nType 'help' for available commands.")
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
		"prop":         "links|extracts",
		"pllimit":      "max",
		"plnamespace":  "0",
		"exintro":      "",
		"explaintext":  "",
		"redirects":    "1",
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
	return res
}
func getSpecificArticle(t string) Response { // Function to get a specific article (based on user input) and return a response struct
	fmt.Printf("------------- Getting article for %s -----------------", t)
	baseURL := "https://en.wikipedia.org/w/api.php?"
	var params = map[string]string{
		"action":      "query",
		"format":      "json",
		"prop":        "links|extracts",
		"pllimit":     "max",
		"plnamespace": "0",
		"exintro":     "",
		"explaintext": "",
		"redirects":   "1",
		"titles":      url.QueryEscape(t),
	}
	for k, v := range params {
		var param = fmt.Sprintf("&%s=%s", k, v)
		baseURL += param
	}
	fmt.Println("\nSpecific article URL:", baseURL)
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
func checkIsLinkValid(linkSlice []Link, choice string) (bool, string) { // Function to check that the user has inputted a valid link (based on the current page)
	for _, v := range linkSlice {
		if strings.EqualFold(v.Title, choice) {
			return true, v.Title
		}
	}
	return false, ""
}
func getAndDisplayLinks(p Page) []Link { // Function that prints all the links in a page to the console, in a easier-to-read format.
	fmt.Printf("\nLinks from the page `%s`:\n", p.Title)
	for i, v := range p.Links {
		if i%6 != 0 || i == 0 {
			fmt.Printf("%-30s ", v.Title)
		} else {
			fmt.Printf("%-30s\n ", v.Title)
		}
	}
	fmt.Println()
	return p.Links
}
func getPageFromResponse(res Response) Page { // Function to extract the page from the response
	var p Page
	for _, v := range res.Query.Pages { // Loop is needed since pages are returned in a map of ints to Pages, but the int key is not known beforehand. Only a single page is returned.
		p = v
	}
	return p
}
func startGame(startResponse Response, t string) (int, []string) {
	currentResponse := startResponse
	startPage := getPageFromResponse(startResponse)
	count := 0
	var clickedLinks = []string{startPage.Title}
	currentPage := startPage
	for currentPage.Title != t {
		currentLinks := getAndDisplayLinks(currentPage)
		if len(currentLinks) == 500 {
			currentLinks = getAdditionalLinks(currentResponse, currentPage.Title, currentLinks)
		}
		userChoice := getUserInput()
		isLinkValid, validChoice := checkIsLinkValid(currentLinks, userChoice)
		if isLinkValid {
			currentResponse = getSpecificArticle(validChoice)
			currentPage = getPageFromResponse(currentResponse)
			clickedLinks = append(clickedLinks, currentPage.Title)
			fmt.Printf("Target page: %s\nCurrent page:\n%s\n%s\n", t, currentPage.Title, currentPage.Extract)
			count++
		} else {
			fmt.Println("Link not found !")
		}
	}
	return count, clickedLinks
}
func getAdditionalLinks(r Response, t string, firstLinks []Link) []Link {
	baseURL := "https://en.wikipedia.org/w/api.php?"
	done := false
	var params = map[string]string{
		"action":      "query",
		"format":      "json",
		"prop":        "links|extracts",
		"pllimit":     "max",
		"plnamespace": "0",
		"exintro":     "",
		"explaintext": "",
		"redirects":   "1",
		"titles":      url.QueryEscape(t),
	}
	for k, v := range params {
		var param = fmt.Sprintf("&%s=%s", k, v)
		baseURL += param
	}
	for !done {
		baseURL2 := "https://en.wikipedia.org/w/api.php?"
		var params = map[string]string{
			"action":      "query",
			"format":      "json",
			"prop":        "links|extracts",
			"pllimit":     "max",
			"plnamespace": "0",
			"exintro":     "",
			"explaintext": "",
			"redirects":   "1",
			"titles":      url.QueryEscape(t),
			"continue":    url.QueryEscape(r.Continue.Continue),
			"plcontinue":  url.QueryEscape(r.Continue.Plcontinue),
		}
		for k, v := range params {
			var param = fmt.Sprintf("&%s=%s", k, v)
			baseURL2 += param
		}
		fmt.Println()
		res, err := http.Get(baseURL2)
		if err != nil {
			panic(err)
		}
		byteSlice, err := io.ReadAll(res.Body)
		if err != nil {
			panic(err)
		}
		r = processResponse(byteSlice)
		nextPage := getPageFromResponse(r)
		firstLinks = append(firstLinks, getAndDisplayLinks(nextPage)...)
		if r.Continue.Continue == "" && r.Continue.Plcontinue == "" {
			fmt.Printf("---------- End of links ----------  %v\n", r.BatchComplete)
			done = true
		}
	}
	return firstLinks
}

/*
upon typing getRandom- get a random article - set it's title to be the target page
then, upon typing start, get another random article (the start page)

 get a list of links from the current page (and display them)
 if the response for the requested page has a continue struct, then make another request with the values in the continue struct, and append the links from that request to the array of links.
 get user input (of a link that they want to "click" on)
 check that link is valid
 fetch that link
 set that link to the current page
 get a list of valid links from the current page...


 ...
 check that link is valid
 if that links is also the target page, end the game and print the count, (and time.Now?)
*/

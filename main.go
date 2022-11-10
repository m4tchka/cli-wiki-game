package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	userInput := getUserInput()
	fmt.Println(userInput)

}

func getUserInput() string {
	r := bufio.NewReader(os.Stdin)
	line, err := r.ReadString('\n')
	if err != nil {
		panic(err)
	}
	return strings.ToLower(strings.TrimSpace(line))
}

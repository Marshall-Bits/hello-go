// Package main is the entry point of the program.
package main

// Importing the fmt package, this package is used for printing to the console.
import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	// Print to the console.
	fmt.Println("Hello, World!")

	// fetch data
	resp, err := http.Get("https://jsonplaceholder.typicode.com/posts")

	if err != nil {
		fmt.Println("Error fetching data")
	}

	defer resp.Body.Close() // if we don't close the body, it will cause a memory leak

	// read the response body
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		fmt.Println("Error reading data")
	}

	fmt.Println(string(body))
}


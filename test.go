package main

import (
	"fmt"
	"strings"
)

// Sample Go code to test the vx editor features
func main() {
	// Test syntax highlighting
	message := "Hello, World!"
	fmt.Println(message)
	
	// Test line wrapping with long lines
	longLine := "This is a very long line that should wrap around when displayed in the editor to test the line wrapping functionality and make sure everything works correctly with mouse selection and scrolling"
	fmt.Println(longLine)
	
	// Test search functionality
	fruits := []string{"find&replace", "banana", "cherry", "date", "elderberry"}
	for i, fruit := range fruits {
		fmt.Printf("%d: %s\n", i, fruit)
	}
	
	// Test multiple buffers
	numbers := make([]int, 10)
	for i := range numbers {
		numbers[i] = i * 2
	}
	
	// Test undo/redo - try editing this
	result := calculate(5, 3)
	fmt.Println("Result:", result)
	
	// Test mouse selection - try selecting this text
	poem := `
Roses are red,
Violets are blue,
VX is awesome,
And so are you!
`
	fmt.Println(strings.TrimSpace(poem))
}

func calculate(a, b int) int {
	return a + b
}

// Test find and replace - try replacing "test" with "demo"
func testFunction() {
	test := "This is a test"
	fmt.Println(test)
	// Another test here
	testValue := 42
	fmt.Println("Test value:", testValue)
}

// Test copy/paste functionality
type Person struct {
	Name string
	Age  int
	City string
}

func createPerson(name string, age int, city string) Person {
	return Person{
		Name: name,
		Age:  age,
		City: city,
	}
}
package main

import (
	"fmt"
	"strings"
	"time"
)

// This is a large test file to stress test the vx editor performance
// It contains ~1000 lines with various line lengths to test:
// - Scrolling performance
// - Line wrapping with long lines
// - Syntax highlighting
// - Search performance
// - Rendering speed

func main() {
	fmt.Println("VX Editor Performance Test File")
	fmt.Println("================================")
	
	// Test with various data structures
	testArrays()
	testMaps()
	testStructs()
	testInterfaces()
	testChannels()
	testGoroutines()
	
	// Long line test - this line should wrap around multiple times when displayed in the editor and test the line wrapping functionality to ensure it handles very long lines correctly without performance degradation or rendering issues
	longLineTest := "This is an extremely long line that contains a lot of text to test how the editor handles line wrapping and rendering performance when dealing with lines that exceed the terminal width by a significant margin and need to be wrapped across multiple visual rows in the display"
	fmt.Println(longLineTest)
}

func testArrays() {
	// Array operations
	numbers := [10]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	for i, num := range numbers {
		fmt.Printf("Index %d: %d\n", i, num)
	}
	
	// Slice operations
	slice := []string{"apple", "banana", "cherry", "date", "elderberry", "fig", "grape", "honeydew"}
	for _, fruit := range slice {
		fmt.Println(fruit)
	}
	
	// Multi-dimensional arrays
	matrix := [][]int{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	}
	
	for i, row := range matrix {
		for j, val := range row {
			fmt.Printf("matrix[%d][%d] = %d\n", i, j, val)
		}
	}
}

func testMaps() {
	// Map operations
	ages := map[string]int{
		"Alice":   25,
		"Bob":     30,
		"Charlie": 35,
		"David":   40,
		"Eve":     45,
	}
	
	for name, age := range ages {
		fmt.Printf("%s is %d years old\n", name, age)
	}
	
	// Nested maps
	users := map[string]map[string]interface{}{
		"user1": {
			"name":  "John Doe",
			"email": "john@example.com",
			"age":   28,
		},
		"user2": {
			"name":  "Jane Smith",
			"email": "jane@example.com",
			"age":   32,
		},
	}
	
	for id, user := range users {
		fmt.Printf("User ID: %s\n", id)
		for key, value := range user {
			fmt.Printf("  %s: %v\n", key, value)
		}
	}
}

func testStructs() {
	type Person struct {
		Name    string
		Age     int
		Email   string
		Address Address
	}
	
	type Address struct {
		Street  string
		City    string
		Country string
		ZipCode string
	}
	
	people := []Person{
		{
			Name:  "Alice Johnson",
			Age:   28,
			Email: "alice@example.com",
			Address: Address{
				Street:  "123 Main St",
				City:    "New York",
				Country: "USA",
				ZipCode: "10001",
			},
		},
		{
			Name:  "Bob Williams",
			Age:   35,
			Email: "bob@example.com",
			Address: Address{
				Street:  "456 Oak Ave",
				City:    "Los Angeles",
				Country: "USA",
				ZipCode: "90001",
			},
		},
		{
			Name:  "Charlie Brown",
			Age:   42,
			Email: "charlie@example.com",
			Address: Address{
				Street:  "789 Pine Rd",
				City:    "Chicago",
				Country: "USA",
				ZipCode: "60601",
			},
		},
	}
	
	for _, person := range people {
		fmt.Printf("Name: %s\n", person.Name)
		fmt.Printf("Age: %d\n", person.Age)
		fmt.Printf("Email: %s\n", person.Email)
		fmt.Printf("Address: %s, %s, %s %s\n",
			person.Address.Street,
			person.Address.City,
			person.Address.Country,
			person.Address.ZipCode)
		fmt.Println("---")
	}
}

func testInterfaces() {
	type Shape interface {
		Area() float64
		Perimeter() float64
	}
	
	type Rectangle struct {
		Width  float64
		Height float64
	}
	
	func (r Rectangle) Area() float64 {
		return r.Width * r.Height
	}
	
	func (r Rectangle) Perimeter() float64 {
		return 2 * (r.Width + r.Height)
	}
	
	type Circle struct {
		Radius float64
	}
	
	func (c Circle) Area() float64 {
		return 3.14159 * c.Radius * c.Radius
	}
	
	func (c Circle) Perimeter() float64 {
		return 2 * 3.14159 * c.Radius
	}
	
	shapes := []Shape{
		Rectangle{Width: 10, Height: 5},
		Circle{Radius: 7},
		Rectangle{Width: 8, Height: 12},
		Circle{Radius: 3.5},
	}
	
	for i, shape := range shapes {
		fmt.Printf("Shape %d:\n", i+1)
		fmt.Printf("  Area: %.2f\n", shape.Area())
		fmt.Printf("  Perimeter: %.2f\n", shape.Perimeter())
	}
}

func testChannels() {
	// Buffered channel
	ch := make(chan int, 5)
	
	go func() {
		for i := 0; i < 5; i++ {
			ch <- i
			fmt.Printf("Sent: %d\n", i)
		}
		close(ch)
	}()
	
	for val := range ch {
		fmt.Printf("Received: %d\n", val)
	}
	
	// Multiple channels with select
	ch1 := make(chan string)
	ch2 := make(chan string)
	
	go func() {
		time.Sleep(1 * time.Second)
		ch1 <- "Message from channel 1"
	}()
	
	go func() {
		time.Sleep(2 * time.Second)
		ch2 <- "Message from channel 2"
	}()
	
	for i := 0; i < 2; i++ {
		select {
		case msg1 := <-ch1:
			fmt.Println(msg1)
		case msg2 := <-ch2:
			fmt.Println(msg2)
		}
	}
}

func testGoroutines() {
	// Worker pool pattern
	jobs := make(chan int, 100)
	results := make(chan int, 100)
	
	// Start workers
	for w := 1; w <= 3; w++ {
		go worker(w, jobs, results)
	}
	
	// Send jobs
	for j := 1; j <= 9; j++ {
		jobs <- j
	}
	close(jobs)
	
	// Collect results
	for a := 1; a <= 9; a++ {
		<-results
	}
}

func worker(id int, jobs <-chan int, results chan<- int) {
	for j := range jobs {
		fmt.Printf("Worker %d processing job %d\n", id, j)
		time.Sleep(time.Second)
		results <- j * 2
	}
}

// Additional test functions to increase file size
func testStringOperations() {
	text := "The quick brown fox jumps over the lazy dog"
	
	// String manipulation
	upper := strings.ToUpper(text)
	lower := strings.ToLower(text)
	title := strings.Title(text)
	
	fmt.Println("Original:", text)
	fmt.Println("Upper:", upper)
	fmt.Println("Lower:", lower)
	fmt.Println("Title:", title)
	
	// String splitting and joining
	words := strings.Split(text, " ")
	for i, word := range words {
		fmt.Printf("Word %d: %s\n", i+1, word)
	}
	
	joined := strings.Join(words, "-")
	fmt.Println("Joined:", joined)
	
	// String searching
	contains := strings.Contains(text, "fox")
	index := strings.Index(text, "fox")
	count := strings.Count(text, "o")
	
	fmt.Printf("Contains 'fox': %v\n", contains)
	fmt.Printf("Index of 'fox': %d\n", index)
	fmt.Printf("Count of 'o': %d\n", count)
}

func testErrorHandling() {
	// Error handling patterns
	result, err := divide(10, 2)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Result: %f\n", result)
	}
	
	result, err = divide(10, 0)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Result: %f\n", result)
	}
	
	// Custom errors
	type CustomError struct {
		Code    int
		Message string
	}
	
	func (e *CustomError) Error() string {
		return fmt.Sprintf("Error %d: %s", e.Code, e.Message)
	}
	
	err = &CustomError{Code: 404, Message: "Not Found"}
	fmt.Println(err)
}

func divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, fmt.Errorf("division by zero")
	}
	return a / b, nil
}

func testPointers() {
	// Pointer basics
	x := 42
	p := &x
	
	fmt.Printf("Value of x: %d\n", x)
	fmt.Printf("Address of x: %p\n", p)
	fmt.Printf("Value at address: %d\n", *p)
	
	*p = 100
	fmt.Printf("New value of x: %d\n", x)
	
	// Pointer to struct
	type Point struct {
		X, Y int
	}
	
	point := Point{X: 10, Y: 20}
	pointPtr := &point
	
	fmt.Printf("Point: (%d, %d)\n", point.X, point.Y)
	pointPtr.X = 30
	fmt.Printf("Modified Point: (%d, %d)\n", point.X, point.Y)
}

func testDefer() {
	// Defer statements
	defer fmt.Println("This is printed last")
	defer fmt.Println("This is printed second to last")
	defer fmt.Println("This is printed third to last")
	
	fmt.Println("This is printed first")
	
	// Defer with function calls
	for i := 0; i < 5; i++ {
		defer fmt.Printf("Deferred: %d\n", i)
	}
	
	fmt.Println("Loop completed")
}

func testPanic() {
	// Panic and recover
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recovered from panic: %v\n", r)
		}
	}()
	
	fmt.Println("About to panic")
	panic("Something went wrong!")
	fmt.Println("This will not be printed")
}

// More test data to reach 1000 lines
var loremIpsum = `
Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.

Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.

Sed ut perspiciatis unde omnis iste natus error sit voluptatem accusantium doloremque laudantium, totam rem aperiam, eaque ipsa quae ab illo inventore veritatis et quasi architecto beatae vitae dicta sunt explicabo.

Nemo enim ipsam voluptatem quia voluptas sit aspernatur aut odit aut fugit, sed quia consequuntur magni dolores eos qui ratione voluptatem sequi nesciunt. Neque porro quisquam est, qui dolorem ipsum quia dolor sit amet, consectetur, adipisci velit.

At vero eos et accusamus et iusto odio dignissimos ducimus qui blanditiis praesentium voluptatum deleniti atque corrupti quos dolores et quas molestias excepturi sint occaecati cupiditate non provident, similique sunt in culpa qui officia deserunt mollitia animi, id est laborum et dolorum fuga.

Et harum quidem rerum facilis est et expedita distinctio. Nam libero tempore, cum soluta nobis est eligendi optio cumque nihil impedit quo minus id quod maxime placeat facere possimus, omnis voluptas assumenda est, omnis dolor repellendus.

Temporibus autem quibusdam et aut officiis debitis aut rerum necessitatibus saepe eveniet ut et voluptates repudiandae sint et molestiae non recusandae. Itaque earum rerum hic tenetur a sapiente delectus, ut aut reiciendis voluptatibus maiores alias consequatur aut perferendis doloribus asperiores repellat.
`

// Test data structures
type User struct {
	ID        int
	Username  string
	Email     string
	FirstName string
	LastName  string
	Age       int
	Active    bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Product struct {
	ID          int
	Name        string
	Description string
	Price       float64
	Stock       int
	Category    string
	Tags        []string
	Rating      float64
	Reviews     []Review
}

type Review struct {
	ID        int
	UserID    int
	Rating    int
	Comment   string
	CreatedAt time.Time
}

type Order struct {
	ID         int
	UserID     int
	Products   []OrderItem
	Total      float64
	Status     string
	CreatedAt  time.Time
	ShippedAt  *time.Time
	DeliveredAt *time.Time
}

type OrderItem struct {
	ProductID int
	Quantity  int
	Price     float64
}

// Sample data generators
func generateUsers(count int) []User {
	users := make([]User, count)
	for i := 0; i < count; i++ {
		users[i] = User{
			ID:        i + 1,
			Username:  fmt.Sprintf("user%d", i+1),
			Email:     fmt.Sprintf("user%d@example.com", i+1),
			FirstName: fmt.Sprintf("First%d", i+1),
			LastName:  fmt.Sprintf("Last%d", i+1),
			Age:       20 + (i % 50),
			Active:    i%2 == 0,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
	}
	return users
}

func generateProducts(count int) []Product {
	products := make([]Product, count)
	categories := []string{"Electronics", "Clothing", "Books", "Home", "Sports"}
	
	for i := 0; i < count; i++ {
		products[i] = Product{
			ID:          i + 1,
			Name:        fmt.Sprintf("Product %d", i+1),
			Description: fmt.Sprintf("This is a description for product %d", i+1),
			Price:       float64(10 + (i * 5)),
			Stock:       100 + (i * 10),
			Category:    categories[i%len(categories)],
			Tags:        []string{fmt.Sprintf("tag%d", i), fmt.Sprintf("tag%d", i+1)},
			Rating:      3.0 + float64(i%3),
			Reviews:     []Review{},
		}
	}
	return products
}

// Algorithm implementations for testing
func bubbleSort(arr []int) []int {
	n := len(arr)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if arr[j] > arr[j+1] {
				arr[j], arr[j+1] = arr[j+1], arr[j]
			}
		}
	}
	return arr
}

func quickSort(arr []int) []int {
	if len(arr) < 2 {
		return arr
	}
	
	pivot := arr[0]
	var less, greater []int
	
	for _, val := range arr[1:] {
		if val <= pivot {
			less = append(less, val)
		} else {
			greater = append(greater, val)
		}
	}
	
	return append(append(quickSort(less), pivot), quickSort(greater)...)
}

func binarySearch(arr []int, target int) int {
	left, right := 0, len(arr)-1
	
	for left <= right {
		mid := left + (right-left)/2
		
		if arr[mid] == target {
			return mid
		} else if arr[mid] < target {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	
	return -1
}

func fibonacci(n int) int {
	if n <= 1 {
		return n
	}
	return fibonacci(n-1) + fibonacci(n-2)
}

func factorial(n int) int {
	if n <= 1 {
		return 1
	}
	return n * factorial(n-1)
}

// Data structure implementations
type Stack struct {
	items []interface{}
}

func (s *Stack) Push(item interface{}) {
	s.items = append(s.items, item)
}

func (s *Stack) Pop() interface{} {
	if len(s.items) == 0 {
		return nil
	}
	item := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return item
}

func (s *Stack) Peek() interface{} {
	if len(s.items) == 0 {
		return nil
	}
	return s.items[len(s.items)-1]
}

func (s *Stack) IsEmpty() bool {
	return len(s.items) == 0
}

type Queue struct {
	items []interface{}
}

func (q *Queue) Enqueue(item interface{}) {
	q.items = append(q.items, item)
}

func (q *Queue) Dequeue() interface{} {
	if len(q.items) == 0 {
		return nil
	}
	item := q.items[0]
	q.items = q.items[1:]
	return item
}

func (q *Queue) Front() interface{} {
	if len(q.items) == 0 {
		return nil
	}
	return q.items[0]
}

func (q *Queue) IsEmpty() bool {
	return len(q.items) == 0
}

type Node struct {
	Value int
	Next  *Node
}

type LinkedList struct {
	Head *Node
}

func (ll *LinkedList) Insert(value int) {
	newNode := &Node{Value: value}
	if ll.Head == nil {
		ll.Head = newNode
		return
	}
	
	current := ll.Head
	for current.Next != nil {
		current = current.Next
	}
	current.Next = newNode
}

func (ll *LinkedList) Delete(value int) {
	if ll.Head == nil {
		return
	}
	
	if ll.Head.Value == value {
		ll.Head = ll.Head.Next
		return
	}
	
	current := ll.Head
	for current.Next != nil {
		if current.Next.Value == value {
			current.Next = current.Next.Next
			return
		}
		current = current.Next
	}
}

func (ll *LinkedList) Search(value int) bool {
	current := ll.Head
	for current != nil {
		if current.Value == value {
			return true
		}
		current = current.Next
	}
	return false
}

// Tree structures
type TreeNode struct {
	Value int
	Left  *TreeNode
	Right *TreeNode
}

type BinaryTree struct {
	Root *TreeNode
}

func (bt *BinaryTree) Insert(value int) {
	bt.Root = insertNode(bt.Root, value)
}

func insertNode(node *TreeNode, value int) *TreeNode {
	if node == nil {
		return &TreeNode{Value: value}
	}
	
	if value < node.Value {
		node.Left = insertNode(node.Left, value)
	} else {
		node.Right = insertNode(node.Right, value)
	}
	
	return node
}

func (bt *BinaryTree) InOrderTraversal() []int {
	var result []int
	inOrder(bt.Root, &result)
	return result
}

func inOrder(node *TreeNode, result *[]int) {
	if node == nil {
		return
	}
	
	inOrder(node.Left, result)
	*result = append(*result, node.Value)
	inOrder(node.Right, result)
}

func (bt *BinaryTree) PreOrderTraversal() []int {
	var result []int
	preOrder(bt.Root, &result)
	return result
}

func preOrder(node *TreeNode, result *[]int) {
	if node == nil {
		return
	}
	
	*result = append(*result, node.Value)
	preOrder(node.Left, result)
	preOrder(node.Right, result)
}

func (bt *BinaryTree) PostOrderTraversal() []int {
	var result []int
	postOrder(bt.Root, &result)
	return result
}

func postOrder(node *TreeNode, result *[]int) {
	if node == nil {
		return
	}
	
	postOrder(node.Left, result)
	postOrder(node.Right, result)
	*result = append(*result, node.Value)
}

// Graph structures
type Graph struct {
	Vertices map[int][]int
}

func NewGraph() *Graph {
	return &Graph{
		Vertices: make(map[int][]int),
	}
}

func (g *Graph) AddEdge(from, to int) {
	g.Vertices[from] = append(g.Vertices[from], to)
}

func (g *Graph) BFS(start int) []int {
	visited := make(map[int]bool)
	queue := []int{start}
	var result []int
	
	for len(queue) > 0 {
		vertex := queue[0]
		queue = queue[1:]
		
		if !visited[vertex] {
			visited[vertex] = true
			result = append(result, vertex)
			
			for _, neighbor := range g.Vertices[vertex] {
				if !visited[neighbor] {
					queue = append(queue, neighbor)
				}
			}
		}
	}
	
	return result
}

func (g *Graph) DFS(start int) []int {
	visited := make(map[int]bool)
	var result []int
	g.dfsHelper(start, visited, &result)
	return result
}

func (g *Graph) dfsHelper(vertex int, visited map[int]bool, result *[]int) {
	visited[vertex] = true
	*result = append(*result, vertex)
	
	for _, neighbor := range g.Vertices[vertex] {
		if !visited[neighbor] {
			g.dfsHelper(neighbor, visited, result)
		}
	}
}

// Hash table implementation
type HashTable struct {
	buckets []*bucket
	size    int
}

type bucket struct {
	items []item
}

type item struct {
	key   string
	value interface{}
}

func NewHashTable(size int) *HashTable {
	buckets := make([]*bucket, size)
	for i := range buckets {
		buckets[i] = &bucket{}
	}
	return &HashTable{
		buckets: buckets,
		size:    size,
	}
}

func (ht *HashTable) hash(key string) int {
	hash := 0
	for _, char := range key {
		hash += int(char)
	}
	return hash % ht.size
}

func (ht *HashTable) Set(key string, value interface{}) {
	index := ht.hash(key)
	bucket := ht.buckets[index]
	
	for i, item := range bucket.items {
		if item.key == key {
			bucket.items[i].value = value
			return
		}
	}
	
	bucket.items = append(bucket.items, item{key: key, value: value})
}

func (ht *HashTable) Get(key string) (interface{}, bool) {
	index := ht.hash(key)
	bucket := ht.buckets[index]
	
	for _, item := range bucket.items {
		if item.key == key {
			return item.value, true
		}
	}
	
	return nil, false
}

func (ht *HashTable) Delete(key string) {
	index := ht.hash(key)
	bucket := ht.buckets[index]
	
	for i, item := range bucket.items {
		if item.key == key {
			bucket.items = append(bucket.items[:i], bucket.items[i+1:]...)
			return
		}
	}
}

// Additional utility functions
func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func isPalindrome(s string) bool {
	s = strings.ToLower(s)
	left, right := 0, len(s)-1
	
	for left < right {
		if s[left] != s[right] {
			return false
		}
		left++
		right--
	}
	
	return true
}

func findMax(arr []int) int {
	if len(arr) == 0 {
		return 0
	}
	
	max := arr[0]
	for _, val := range arr[1:] {
		if val > max {
			max = val
		}
	}
	
	return max
}

func findMin(arr []int) int {
	if len(arr) == 0 {
		return 0
	}
	
	min := arr[0]
	for _, val := range arr[1:] {
		if val < min {
			min = val
		}
	}
	
	return min
}

func sum(arr []int) int {
	total := 0
	for _, val := range arr {
		total += val
	}
	return total
}

func average(arr []int) float64 {
	if len(arr) == 0 {
		return 0
	}
	return float64(sum(arr)) / float64(len(arr))
}

// More test content to reach 1000 lines
func testConcurrency() {
	// WaitGroup example
	var wg sync.WaitGroup
	
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			fmt.Printf("Goroutine %d starting\n", id)
			time.Sleep(time.Second)
			fmt.Printf("Goroutine %d done\n", id)
		}(i)
	}
	
	wg.Wait()
	fmt.Println("All goroutines completed")
}

func testMutex() {
	// Mutex example
	var mu sync.Mutex
	counter := 0
	
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			counter++
			mu.Unlock()
		}()
	}
	
	wg.Wait()
	fmt.Printf("Final counter value: %d\n", counter)
}

// End of test file - this should be around 1000 lines

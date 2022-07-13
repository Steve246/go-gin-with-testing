package main

import (
	"enigmacamp.com/golatihanlagi/delivery"
	_ "github.com/lib/pq"
)

func main() {
	delivery.NewServer().Run()
}


/*

Unit Testing
masin-masing folder di testing
1. Usecase
2. Unit Testing
3. Controller
*/


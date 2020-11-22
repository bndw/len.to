package main

import (
	"fmt"
)

func main() {
	fmt.Printf(`name: string @index(fulltext, term) .
location: string @index(fulltext, term) .
url: string @index(term) .
tagged: [uid] @count @reverse .
`)
}

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	input := bufio.NewScanner(os.Stdin)

	for input.Scan() {
		fmt.Println(input.Text())
	}

	if err := input.Err(); err != nil {
		log.Fatal(err)
	}
}

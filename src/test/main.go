package main

import (
	"fmt"

	"time"
)

func main() {
	t := time.Now().String()

	date := t[0:10]

	fmt.Println(date)

}

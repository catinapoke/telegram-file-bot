package main

import (
	"fmt"

	"github.com/catinapoke/go-microservice/fileservice"
)

func main() {
	fmt.Println("Hello world!")
	controller := fileservice.CreateController("fileservice", "3001")
	_, err := controller.Get(0)

	if err != nil {
		fmt.Printf("Can't get zero file!")
	}
}

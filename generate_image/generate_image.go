package main

import (
	"emojimage/pkg"
	"fmt"
	"log"
	"os"
	"strconv"
)

func main() {
	path := ""
	size := 20

	switch len(os.Args) {
	case 0:
		fallthrough
	case 1:
		fmt.Println("Usage: go run ./generate_image <path to PNG file> [<width/height>]")
		os.Exit(1)
	case 2:
		path = os.Args[1]
	case 3:
		path = os.Args[1]

		_size, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Invalid width/height")
			os.Exit(1)
		}

		size = _size
	}

	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	image, err := pkg.GenerateImage(f, uint(size))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(image)
}

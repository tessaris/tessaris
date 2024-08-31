package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	data, err := os.ReadFile("./.version")

	if err != nil {
		fmt.Println("Error reading file:", err)
		os.Exit(1)
	}

	s := strings.Split(string(data), ".")
	major, minor, patch := s[0], s[1], s[2]

	// Convert the versions to integers.
	majorInt, err := strconv.Atoi(major)
	if err != nil {
		fmt.Println("Error converting major to int:", err)
		os.Exit(2)
	}

	minorInt, err := strconv.Atoi(minor)
	if err != nil {
		fmt.Println("Error converting minor to int:", err)
		os.Exit(2)
	}

	patchInt, err := strconv.Atoi(patch)
	if err != nil {
		fmt.Println("Error converting patch to int:", err)
		os.Exit(2)
	}

	if os.Args[1] == "major" {
		majorInt++
		minorInt = 0
		patchInt = 0
	} else if os.Args[1] == "minor" {
		minorInt++
		patchInt = 0
	} else if os.Args[1] == "patch" {
		patchInt++
	} else {
		os.Exit(3)
	}

	version := fmt.Sprintf("%d.%d.%d", majorInt, minorInt, patchInt)

	os.WriteFile(".version", []byte(version), 0644)
}

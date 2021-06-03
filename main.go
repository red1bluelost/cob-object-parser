package main

import (
	"os"
)

func main() {
	file, err := os.Open("../spmm.o")
	if err != nil {
		panic(err)
	}
	_, err = ReadBundleObject()
}

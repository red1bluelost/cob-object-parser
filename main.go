package main

import (
	"fmt"
	"os"

	cob "github.com/red1bluelost/cob-object-parser/clangoffloadbundle"
)

func main() {
	file, err := os.Open("../spmm.o")
	if err != nil {
		panic(err)
	}
	_, err = cob.ReadBundleObject(file)
	fmt.Printf("Result: %s\n", err)
}

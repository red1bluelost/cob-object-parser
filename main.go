package main

import (
	"fmt"
	"io/ioutil"
	"os"

	cob "github.com/red1bluelost/cob-object-parser/clangoffloadbundle"
)

func main() {
	file, err := os.Open("../spmm.o")
	if err != nil {
		panic(err)
	}
	obj, err := cob.ReadBundleObject(file)
	fmt.Printf("Result: %s\nError: %s\n", obj, err)

	err = ioutil.WriteFile("../spmm.hsaco", obj.GetCodeObject(1), 0666)
}

package clangoffloadbundle

import (
	"bufio"
	"fmt"
	"io"
)

type ObjectLayout struct {
	numBundleEntries uint64
	headers          []bundleEntryHeader
	codeObjects      [][]byte
}

type bundleEntryHeader struct {
	offset uint64
	size   uint64
	idLen  uint64
	id     []byte
}

func ReadBundleObject(f io.Reader) (*ObjectLayout, error) {
	inFile := bufio.NewReader(f)
	magicString := make([]byte, 24)
	if n, err := inFile.Read(magicString); err != nil || n != 24 {
		return nil, fmt.Errorf("read n=%d, and error: %s", n, err)
	}
	return nil, fmt.Errorf("success: %s", magicString)
}

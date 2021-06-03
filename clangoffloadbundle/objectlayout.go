package clangoffloadbundle

import (
	"bufio"
	"fmt"
	"io"
	"reflect"
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

func ReadBundleObject(file io.Reader) (*ObjectLayout, error) {
	r := bufio.NewReader(file)
	if err := verifyMagicString(r); err != nil {
		return nil, err
	}

	return nil, fmt.Errorf("to be implemented")
}

func verifyMagicString(r *bufio.Reader) error {
	magicString := make([]byte, 24)
	if n, err := r.Read(magicString); err != nil || n != 24 {
		return fmt.Errorf("read n=%d, and error: %s", n, err)
	}
	if !reflect.DeepEqual(magicString, []byte("__CLANG_OFFLOAD_BUNDLE__")) {
		return fmt.Errorf("magic string not located at front of file")
	}
	return nil
}
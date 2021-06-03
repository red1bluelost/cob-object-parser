package clangoffloadbundle

import (
	"bufio"
	"encoding/binary"
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
	var err error
	if err = verifyMagicString(r); err != nil {
		return nil, err
	}
	objLayout := new(ObjectLayout)

	if objLayout.numBundleEntries, err = readNumber(r); err != nil {
		return nil, err
	}

	return objLayout, fmt.Errorf("to be implemented")
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

func readNumber(r *bufio.Reader) (uint64, error) {
	number8byte := make([]byte, 8)
	if n, err := r.Read(number8byte); err != nil || n != 8 {
		return 0, fmt.Errorf("read n=%d, and error: %s", n, err)
	}
	return binary.LittleEndian.Uint64(number8byte), nil
}

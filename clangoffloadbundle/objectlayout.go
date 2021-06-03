package clangoffloadbundle

import (
	"bufio"
	"bytes"
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

func (o *ObjectLayout) String() string {
	var b bytes.Buffer
	b.WriteString(fmt.Sprintf("{%d [", o.numBundleEntries))
	for i, header := range o.headers {
		b.WriteString(fmt.Sprintf("%d:%s", i, header.String()))
		if uint64(i) != o.numBundleEntries-1 {
			b.WriteString(" ")
		}
	}
	b.WriteString("] [")
	for i, object := range o.codeObjects {
		b.WriteString(fmt.Sprintf("%d:len=%d", i, len(object)))
		if uint64(i) != o.numBundleEntries-1 {
			b.WriteString(" ")
		}
	}
	b.WriteString("]}")
	return b.String()
}

type bundleEntryHeader struct {
	offset uint64
	size   uint64
	idLen  uint64
	id     []byte
}

func (b *bundleEntryHeader) String() string {
	return fmt.Sprintf("{%d %d %d %s}", b.offset, b.size, b.idLen, b.id)
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

	objLayout.headers = make([]bundleEntryHeader, 0)
	for i := uint64(0); i < objLayout.numBundleEntries; i++ {
		header, err := readHeader(r)
		if err != nil {
			return nil, err
		}
		objLayout.headers = append(objLayout.headers, header)
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

func readHeader(r *bufio.Reader) (bundleEntryHeader, error) {
	header := bundleEntryHeader{}
	var err error
	header.offset, err = readNumber(r)
	if err != nil {
		return bundleEntryHeader{}, err
	}
	header.size, err = readNumber(r)
	if err != nil {
		return bundleEntryHeader{}, err
	}
	header.idLen, err = readNumber(r)
	if err != nil {
		return bundleEntryHeader{}, err
	}
	idBytes := make([]byte, header.idLen)
	if n, err := r.Read(idBytes); err != nil || uint64(n) != header.idLen {
		return bundleEntryHeader{}, fmt.Errorf("read n=%d, and error: %s", n, err)
	}
	header.id = idBytes
	return header, err
}

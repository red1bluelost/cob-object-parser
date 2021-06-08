package clangoffloadbundle

import (
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

	readSeeker io.ReadSeeker
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

func (o *ObjectLayout) GetCodeObject(index int) []byte {
	if o.codeObjects[index] == nil {
		header := o.headers[index]
		r := o.readSeeker
		if _, err := r.Seek(int64(header.offset), io.SeekStart); err != nil {
			panic(err)
		}
		object := make([]byte, header.size)
		if n, err := io.ReadFull(r, object); err != nil || uint64(n) != header.size {
			panic(fmt.Errorf("read %d but should be %d, and error: %s", n, header.size, err))
		}
		o.codeObjects[index] = object
	}
	return o.codeObjects[index]
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

func ReadBundleObject(r io.ReadSeeker) (*ObjectLayout, error) {
	var err error
	if err = verifyMagicString(r); err != nil {
		return nil, err
	}
	objLayout := new(ObjectLayout)
	objLayout.readSeeker = r

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

	objLayout.codeObjects = make([][]byte, objLayout.numBundleEntries)
	return objLayout, nil
}

func verifyMagicString(r io.Reader) error {
	magicString := make([]byte, 24)
	if n, err := r.Read(magicString); err != nil || n != 24 {
		return fmt.Errorf("read n=%d, and error: %s", n, err)
	}
	if !reflect.DeepEqual(magicString, []byte("__CLANG_OFFLOAD_BUNDLE__")) {
		return fmt.Errorf("magic string not located at front of file")
	}
	return nil
}

func readNumber(r io.Reader) (uint64, error) {
	number8byte := make([]byte, 8)
	if n, err := r.Read(number8byte); err != nil || n != 8 {
		return 0, fmt.Errorf("read n=%d, and error: %s", n, err)
	}
	return binary.LittleEndian.Uint64(number8byte), nil
}

func readHeader(r io.Reader) (bundleEntryHeader, error) {
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

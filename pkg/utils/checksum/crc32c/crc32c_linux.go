// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package crc32 implements the 32-bit cyclic redundancy check, or CRC-32,
// checksum. See http://en.wikipedia.org/wiki/Cyclic_redundancy_check for
// information.
package crc32c

import (
	"hash"
	"io"
)

// The size of a CRC-32 checksum in bytes.
const Size = 4

// digest represents the partial evaluation of a checksum.
type digest struct {
	crc uint32
}

// New creates a new hash.Hash32 computing the CRC-32 checksum
// using the polynomial represented by the Table.
func New() hash.Hash32 {
	return &digest{crc: 0}
}

func (d *digest) Size() int { return Size }

func (d *digest) BlockSize() int { return 1 }

func (d *digest) Sum(in []byte) []byte {
	s := d.crc
	return append(in, byte(s>>24), byte(s>>16), byte(s>>8), byte(s))
}

func (d *digest) Sum32() uint32 { return d.crc }

func (d *digest) Reset() { d.crc = 0 }

// Update returns the result of adding the bytes in p to the crc.
func (d *digest) update(crc uint32, p []byte) uint32 {
	return updateCastanagoliPCL(crc, p)
}

func (d *digest) Write(p []byte) (n int, err error) {
	d.crc = d.update(d.crc, p)
	return len(p), nil
}

// Convenience functions

func Sum32(data []byte) uint32 {
	crc32 := New()
	crc32.Reset()
	crc32.Write(data)
	return crc32.Sum32()
}

func Sum(reader io.Reader) (uint32, error) {
	h := New()
	var err error
	for err == nil {
		length := 0
		byteBuffer := make([]byte, 1024*1024)
		length, err = reader.Read(byteBuffer)
		byteBuffer = byteBuffer[0:length]
		h.Write(byteBuffer)
	}
	if err != io.EOF {
		return 0, err
	}
	return h.Sum32(), nil
}
